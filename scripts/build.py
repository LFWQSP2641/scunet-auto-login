#!/usr/bin/env python3
from __future__ import annotations

import logging
import os
import platform
import subprocess
import sys
from pathlib import Path
from typing import Iterable, List, Optional, Sequence

logging.basicConfig(level=logging.INFO, format="%(levelname)s: %(message)s")
logger = logging.getLogger(__name__)

ANDROID_LICENSE_REL = Path("licenses") / "android-sdk-license"
FIXED_NDK_VERSION = "28.0.13004108"
NDK_VERSION_FILE = "source.properties"


def find_sdk() -> tuple[Path, Path]:
    """Locate the Android SDK and NDK, updating relevant environment variables."""
    sdk_path = _find_android_sdk()
    if not sdk_path:
        logger.error("android SDK not found")
        sys.exit(1)

    ndk_path = _find_android_ndk(sdk_path)
    if not ndk_path:
        logger.error("android NDK not found")
        sys.exit(1)

    os.environ["ANDROID_HOME"] = str(sdk_path)
    os.environ["ANDROID_SDK_HOME"] = str(sdk_path)
    os.environ["ANDROID_NDK_HOME"] = str(ndk_path)
    os.environ["NDK"] = str(ndk_path)

    llvm_bin = (
        ndk_path / "toolchains" / "llvm" / "prebuilt" / _toolchain_host_dir() / "bin"
    )
    if llvm_bin.exists():
        _prepend_path(llvm_bin)
    else:
        logger.debug("LLVM toolchain bin directory not found at %s", llvm_bin)

    return sdk_path, ndk_path


def _find_android_sdk() -> Optional[Path]:
    search_paths = [
        "$ANDROID_HOME",
        "$HOME/Android/Sdk",
        "$HOME/.local/lib/android/sdk",
        "$HOME/Library/Android/sdk",
    ]
    for raw_path in search_paths:
        expanded = os.path.expandvars(raw_path)
        expanded = os.path.expanduser(expanded)
        candidate = Path(expanded)
        if candidate.is_dir() and (candidate / ANDROID_LICENSE_REL).is_file():
            return candidate.resolve()
    return None


def _find_android_ndk(sdk_path: Path) -> Optional[Path]:
    fixed_path = sdk_path / "ndk" / FIXED_NDK_VERSION
    if (fixed_path / NDK_VERSION_FILE).is_file():
        return fixed_path.resolve()

    env_path = os.environ.get("ANDROID_NDK_HOME")
    if env_path:
        env_candidate = Path(env_path).expanduser()
        if (env_candidate / NDK_VERSION_FILE).is_file():
            return env_candidate.resolve()

    ndk_root = sdk_path / "ndk"
    if not ndk_root.is_dir():
        return None

    candidates = []
    for entry in ndk_root.iterdir():
        if entry.is_dir() and (entry / NDK_VERSION_FILE).is_file():
            candidates.append(entry)

    if not candidates:
        return None

    candidates.sort(key=lambda path: _version_key(path.name), reverse=True)

    selected = candidates[0].resolve()
    if selected.name != FIXED_NDK_VERSION:
        logger.warning(
            "reproducibility warning: using NDK version %s instead of %s",
            selected.name,
            FIXED_NDK_VERSION,
        )
    return selected


def _version_key(version: str) -> List[int]:
    parts = []
    for part in version.split("."):
        try:
            parts.append(int(part))
        except ValueError:
            parts.append(0)
    return parts


def _toolchain_host_dir() -> str:
    goos_map = {
        "windows": "windows",
        "linux": "linux",
        "darwin": "darwin",
    }
    system = platform.system().lower()
    goos = goos_map.get(system, system)
    return f"{goos}-x86_64"


def _prepend_path(path: Path) -> None:
    path_str = str(path)
    current = os.environ.get("PATH", "")
    segments = [segment for segment in current.split(os.pathsep) if segment]
    if path_str not in segments:
        os.environ["PATH"] = os.pathsep.join([path_str] + segments)


def ensure_gopath() -> None:
    if os.environ.get("GOPATH"):
        return
    gopath = _detect_gopath()
    if gopath:
        os.environ["GOPATH"] = gopath
    else:
        logger.error("unable to determine GOPATH")
        sys.exit(1)


def _detect_gopath() -> Optional[str]:
    try:
        output = subprocess.check_output(["go", "env", "GOPATH"], text=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        output = ""
    gopath = output.strip()
    if gopath:
        return gopath
    return str(Path.home() / "go")


def run_command(args: Sequence[str]) -> int:
    try:
        completed = subprocess.run(args, check=False)
    except FileNotFoundError as exc:
        logger.error("failed to execute %s: %s", args[0], exc)
        return 1
    return completed.returncode


def main(argv: Sequence[str]) -> None:
    if len(argv) < 2:
        logger.error("usage: %s <command> [args...]", Path(argv[0]).name)
        sys.exit(1)

    find_sdk()
    ensure_gopath()

    exit_code = run_command(list(argv[1:]))
    if exit_code != 0:
        sys.exit(exit_code)


if __name__ == "__main__":
    main(sys.argv)

# aurm
Download the build files if you have not already installed it, or if the package is outdated.
https://wiki.archlinux.org/index.php/Arch_User_Repository
# Usage
```
$ sudo pacman -Rsn $(pacman -Qdmq) // If you have not removed the orphan package.
$ aurm
$ cd ${PACKAGE_NAME}
$ makepkg -sirc
```

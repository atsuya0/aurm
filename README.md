# aurm
This is a tool for Arch Linux.  
Arch Linux is an independently developed, x86-64 general-purpose GNU/Linux distribution that strives to provide the latest stable versions of most software by following a rolling-release model.  
https://wiki.archlinux.org/title/Arch_Linux  
(https://wiki.archlinux.jp/index.php/Arch_Linux)

Download the build files of the outdated package.  
https://wiki.archlinux.org/index.php/Arch_User_Repository  
(https://wiki.archlinux.jp/index.php/Arch_User_Repository)

# Usage
```
$ sudo pacman -Rsn $(pacman -Qdmq) // If you have not removed the orphan packages.
$ aurm
Download pkg1
Download pkg2
$ cd pkg1
$ makepkg -sirc
```

# aurm
Download the build files if you have not already installed it, or if the package is outdated.
https://wiki.archlinux.org/index.php/Arch_User_Repository
# Usage
```
$ mkdir -p ~/.local/share/aurm
$ echo 'chromium-widevine' > ~/.local/share/aurm/list.txt
$ aurm
$ cd chromium-widevine
$ makepkg -sirc
```

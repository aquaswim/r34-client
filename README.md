### Cross compile to windows

1. install mingw
2. setup fyne
3. set this env
    ```
   GOOS=windows
   CC=x86_64-w64-mingw32-gcc #this is for fedora, maybe different from your distro
   ```
4. run this fyne package

readmore: https://docs.fyne.io/started/cross-compiling

### Acknowledgement

* Icon: [Feathers icon](https://feathericons.com)
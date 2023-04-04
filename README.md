# Run Original Tetris by Alexey Pajitnov on your Computer

In light of the latest movie [about the history of Tetris][4], I decided to
try and play the original version of the game.  I was surprised to find that
it was not available on any of the major platforms.  I was able to find a
version of the game that runs on the MAME emulator, and get it running on my
Mac, sharing the process here in case it helps anyone else.

The solution is based on the [original instructions][5], and the [SDL MAME
tutorial][6].  The only changes I made are:

- dockerized the build process to ensure a consistent environment;
- used the latest version of MAME (0.252);

## Architecture

```
+---- docker --------------------+  +--- MAME ----+      \o 
| SIMH(PDP-11): localhost:2323 <-+--+-15ИЭ-00-013 | <--   |\
+--------------------------------+  +-------------+       Л
          (backend)                   (frontend)        (you)
```

## Prerequistes

### Docker
Download and install docker from https://www.docker.com/products/docker-desktop/

## Quickstart (easy-run)

WIP


## Quickstart (manual)

### build docker compose image
Ensure that you're in the same directory as the `docker-compose.yml` file.

    docker compose build


### Install MAME

#### Linux
1. If you're on Ubuntu, you can do this with:

       sudo apt install mame

   If you're on Gentoo, you most likely know what to do.

#### Mac 
1. Download the SDL framework from [SDL Releases page][3] (`SDL2-2.26.4.dmg` or
   later)
2. Unpack SDL to `/Library/Frameworks`
3. Get the latest MAME:
   - **Intel:** Get the latest [mame][1] (discovered though comments on [this][2] page)
   - **Apple Silicon:** Get the latest [mame][8]
4. Unpack the zip file to ./mame

#### Windows
1. Get the latest release of [MAME][7]
2. Install MAME onto your system by running the downloaded installer.

## Run
Automatic (macOS or Linux):

    make

Manual:

1. Start the PDP-11 simulator:

       docker compose up -d

   This will start the SimH simulating PDP-11, and expose the serial port on
   `localhost:2323`.  It will wait for connection from the terminal.

   Precise build instructions are in Dockerfile, if you're of an inquisitive
   type.

2. Start the MAME emulator that emulates 15ИЭ-00-013 terminal:

       ./mame/mame ie15 -rompath files/rom -window -rs232 null_modem -bitb socket.localhost:2323

   This will instruct MAME to emulate the 15ИЭ-00-013 terminal, and connect
   it's RS232 "port" to the PDP-11 simulator using bitbanger connection to
   localhost:2323.  The `-window` option will start the emulator in a
   separate window.

3. At the prompt (which looks like a "."), type:

         RUN DL1:TETRIS

   And press Enter.

   The terminal has a non-standard keymap, so to type the colon, ":" you need
   to press `[']` key (on US keyboard, it's to the left of Enter or Return
   key).

There are some other games on the DL1: disk, to view them all, type:

    DIR DL1:*

The "*" character can be entered by pressing `[Shift]+[']` keys.

## Other References
1. https://retrogamesultra.com/2019/02/24/running-the-mame-arcade-emulator-on-mac-os-x/
2. https://wiki.mamedev.org/index.php/SDL_Supported_Platforms

## License
This work is licensed under the Creative Commons Attribution-ShareAlike 4.0

Disk images are mirrored from:
- http://astio.ciotoni.net/tetris/rl0.dsk (PDP system image with RT-11 OS)
- http://astio.ciotoni.net/tetris/games.dsk

as per original instructions.

[1]: https://www.mediafire.com/file/zbo90w4rvfey9gb/mame_252.zip/file
[2]: https://sdlmame.lngn.net/2023/02/23/mame-0-252-released/
[3]: https://github.com/libsdl-org/SDL/releases/
[4]: https://www.imdb.com/title/tt10100484/
[5]: https://lab.dyne.org/OriginalTetrisHowto
[6]: http://bamf2048.github.io/sdl_mame_tut/
[7]: https://github.com/mamedev/mame/releases/download/mame0253/mame0253b_64bit.exe
[8]: https://sdlmame.lngn.net/stable/mame0252-arm64.zip

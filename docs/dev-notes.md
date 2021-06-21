# Dev Notes

There are several ambiguous instructions in CHIP-8
http://johnearnest.github.io/Octo/docs/SuperChip.html#compatibility

This emulator implements the newer or quirks mode of those instructions, namely:

- 8XY6 & 8XYE - Vy is ignored
- FX55 & FX65 - I is not incremented


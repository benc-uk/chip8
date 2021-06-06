# Opcode decoding
See the defacto CHIP-8 reference  
http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3.0

Note. NN is used here in place of KK, the meaning is the same

```text

Where K = the kind or type of opcode

Opcode with three params: x, y, n

{   byte 1  |   byte 2  }
+-----+-----+-----+-----+
|  K  |  X  |  Y  |  N  |
+-----+-----+-----+-----+
(all are nibbles)



Opcode with two params: x, y, n

{   byte 1  |   byte 2  }
+-----+-----+-----+-----+
|  K  |  X  |    NN     |
+-----+-----+-----+-----+
(NN is a byte)



Opcode with one param

{   byte 1  |   byte 2  }
+-----+-----+-----+-----+
|  K  |      NNN        |
+-----+-----+-----+-----+
(NNN is 12 bits)
```
# wire-os

-	This is a Go program which will eventually be able to download any OTA, apply patches to it, and pack it to be able to download onto any target (production bots are not supported. only oskr/dev/whiskey).
-	Right now, there is just a very basic debug shell which `download`s 1.6.0.3331.ota, `mount`s it, and `pack`s it. Patches/modifiers are not implemented yet.
-	This is what will replace what is currently known as WireOS.

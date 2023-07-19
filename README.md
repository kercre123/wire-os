# wire-os

-	This is a Go program which will eventually be able to download any OTA, apply patches to it, and pack it to be able to download onto any target (production bots are not supported. only oskr/dev/whiskey).
-   Current status: Functions for mounting, unmounting, downloading, packing, and patching are fully functional. Though, they are just held together with a debug shell program for now.
-   Components:
    -   Download Manager - downloads an OTA from URL, saves it to disk, keeps track of the URL. If it is a `latest.ota`, it reads the manifest to figure out version string
    -   OTA Handler - mounts, unmounts, packs OTAs
    -   OTA Patcher - contains functions for patching OTAs. Patch examples: custom boot animation, changing server env
-   Command list:

```
(for now) These commands should be run without arguments. They are interactive
download - Download an OTA from URL
mount - Mount an OTA that has been downloaded
unmount - Unmount currently-mounted OTA
patch - Apply wireos patches to mounted OTA
pack - Pack mounted OTA
OTA format example: 1.6.0.3331
Targets: 0 = dev, 1 = whiskey, 2 = oskr, 3 = orange
```

-	This program is what will replace what is currently known as WireOS.

## OTA links:

-	http://ota.global.anki-dev-services.com/vic/master-dev/lo8awreh23498sf/full/
-	https://ota-cdn.anki.com/vic/dev/werih2382df23ij/full/
-	https://assets.digitaldreamlabs.com/vic/ufxTn3XGcVNK2YrF/master/dev/vicos-

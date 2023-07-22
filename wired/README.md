## impl

program must do the following:

-	Run a webserver
-	This server has options
	-	Example options:
		1.	freqchange: max, balanced, reg
		2.	behavior modifiers
			-	would edit the victor behavior tree files
		3.	toggle rainbow lights
		4.	easy "features" interface
-	These will be applied to the filesystem
-	These options need init functions
-	wired will run at boot, as multi-user.target
	-	if modifications not made to filesystem, it will do it
-	modifications must be able to be saved and loaded

ideas of how:

-	create functions which apply the mods, each with a way to reset to default

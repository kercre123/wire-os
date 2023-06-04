# wired

-	Will be the manager of the custom settings in WireOS.
-	Current status:
	-	Modifier system in place
		-	Apply, Remove, Init functions for each modifier, can be modified dynamically
	-	Two modifiers implemented
		-	HigherPerformance
		-	NoSnore
	-	Just a debug shell for now, eventually there will be a web UI with wire-pod CSS
	-	Build system in place
		-	Current golang:buster container is broken for armel, decided to go with a stretch archive image instead
		-	Runs natively in his filesystem as armel, this prorgram doesn't really need the benefits of armhf

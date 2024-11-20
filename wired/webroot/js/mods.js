UpdateAllMods()

async function UpdateAllMods(undata) {
    document.getElementById('restartNeeded').style.display = 'none';
    document.getElementById('showDuringVicRestart').style.display = 'none';

    var data = await GetCurrent('FreqChange');
    var radioButtons = document.getElementsByName("frequency");
    for(var i = 0; i < radioButtons.length; i++){
        if(radioButtons[i].value == data.freq){
            radioButtons[i].checked = true;
            break;
        }
    }

    // data = await GetCurrent('RainbowLights');
    // console.log(data.enabled)
    // radioButtons = document.getElementsByName("rainbowlights");
    // for(var i = 0; i < radioButtons.length; i++){
    //     if(radioButtons[i].value == JSON.stringify(data.enabled)){
    //         radioButtons[i].checked = true;
    //         break;
    //     }
    // }

    // let response = await GetCurrent('BootAnim');
    // let checkbox = document.getElementById('bootAnimDefault');
    // let divUpload = document.getElementById('bootAnimUploadHide');

    // if(response.default == false) {
    //     checkbox.checked = false;
    //     divUpload.style.display = "block";

    //     let img = document.createElement('img');
    //     img.src = `data:image/gif;base64,${response.gifdata}`;
    //     document.getElementById('bootAnimCurrent').innerHTML = "";
    //     document.getElementById('bootAnimCurrent').appendChild(img);
    // } else {
    //     document.getElementById('bootAnimCurrent').innerHTML = "";
    //     checkbox.checked = true;
    //     divUpload.style.display = "none";
    // }
    // bootAnimCheckValidate()
}

async function GetCurrent(mod) {
    let response = await fetch('/api/mods/current/' + mod);
    let data = await response.json();
    return data;
}

function SetModStatus(message) {
    statusMsg = document.createElement("h3")
    statusDiv = document.getElementById('modStatus')
    statusMsg.textContent = message
    statusDiv.innerHTML = ""
    document.getElementById('modStatus').style.display = 'block';
    statusDiv.appendChild(statusMsg)
}

function HideModStatus() {
    document.getElementById('modStatus').style.display = 'none';
}

async function SendJSON(mod, json) {
    document.getElementById('mods').style.display = 'none';
    SetModStatus(mod + " is applying, please wait...")
    let response = await fetch('/api/mods/modify/' + mod, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: json,
    });
    let data = await response.json();
    UpdateAllMods(data)
    if (data.status == "success") {
        document.getElementById('mods').style.display = 'block';
        HideModStatus()
    } else {
        document.getElementById('mods').style.display = 'block';
        HideModStatus()
    }
    return data;
}

async function FreqChange_Submit() {
    let freq = document.querySelector('input[name="frequency"]:checked').value;
    let data = await SendJSON('FreqChange', `{"freq":` + freq + `}`);
    console.log('Success:', data);
    CheckIfRestartNeeded("FreqChange");
}

async function RainbowLights_Submit() {
    let enabled = document.querySelector('input[name="rainbowlights"]:checked').value;
    let data = await SendJSON('RainbowLights', `{"enabled":` + enabled + `}`);
    console.log('Success:', data);
    CheckIfRestartNeeded("RainbowLights");
}

/* <div id="wakeWordStatus">
</div>
<button id="startTraining" onclick="startWakeWordFlow()">
    Train a new wake word
</button>
<button id="wakeWordListen" onclick="doListen()" style="display:none">
    Listen
</button>
<button id="genWakeWord" onclick="genWakeWord()" style="display:none">
    Generate Wake Word
</button>
</div> */

function hide(element) {
    document.getElementById(element).style.display = 'none';
}

function show(element) {
    document.getElementById(element).style.display = 'block';
}

function setWakeStatus(status) {
    document.getElementById("wakeWordStatus").innerHTML = ""
    let stat = document.createElement("p")
    stat.innerHTML = status
    document.getElementById("wakeWordStatus").appendChild(stat)
}

let recIndex = 0

async function startWakeWordFlow() {
    recIndex = 0
    hide("startTraining")
    setWakeStatus("Starting listener... Vector's eyes will go dark.")
    await fetch("/api/mods/wakeword/StartListener")
    show("wakeWordListen")
    hide("startTraining")
    setWakeStatus("Listener started. Press 'Listen', wait for the countdown on Vector's screen, then say your wake word to Vector. Do this again at least two more times, then you will be able to click 'Generate Wake Word'.")
}

async function doListen() {
    setWakeStatus("Listen request sent. Look at Vector's screen!")
    hide("startOver")
    hide("wakeWordListen")
    await fetch("/api/mods/wakeword/Listen")
    recIndex++
    show("wakeWordListen")
    if (recIndex >= 3) {
        setWakeStatus("You now have three or more recordings, which means you can generate a wake word. You can create more recordings if you want better accuracy.<br><br>Recordings made: " + recIndex)
        show("genWakeWord")
    } else {
        setWakeStatus("Press 'Listen', wait for the countdown on Vector's screen, then say your wake word to Vector.<br><br>Recordings made: " + recIndex)
    }
    if (recIndex >= 1) {
        show("startOver")
    }
}

async function genWakeWord() {
    setWakeStatus("Generating wake word...")
    hide("genWakeWord")
    hide("wakeWordListen")
    hide("startOver")
    await fetch("/api/mods/wakeword/GenWakeWord")
    setWakeStatus("Wake word generated and installed. Starting anki programs...")
    await fetch("/api/mods/wakeword/StopListener")
    setWakeStatus("Your custom wake word is now implemented.")
    show("startTraining")
}

async function startWakeWordOver() {
    setWakeStatus("Deleting data...")
    await fetch("/api/mods/wakeword/StartOver")
    hide("genWakeWord")
    hide("startOver")
    recIndex = 0
    setWakeStatus("The recordings have been deleted. Press 'Listen', wait for the countdown on Vector's screen, then say your wake word to Vector.")
}

async function CheckIfRestartNeeded(mod) {
    let response = await fetch('/api/mods/needsrestart/' + mod, {
        method: 'POST',
    });
    let data = await response.text()
    if (data.includes("true")) {
        document.getElementById('restartNeeded').style.display = 'block';
    }
}

async function RestartVic() { 
    SetModStatus("")
    document.getElementById("restartButton").disabled = true
    document.getElementById('showDuringVicRestart').style.display = 'block';
    document.getElementById('mods').style.display = 'none';
    fetch('/api/restartvic', {
        method: 'POST',
    }).then(response => {console.log(response); document.getElementById("restartButton").disabled = false; document.getElementById('restartNeeded').style.display = 'none'; document.getElementById('showDuringVicRestart').style.display = 'none'; document.getElementById('mods').style.display = 'block';})
}

async function BootAnim_Test() {
    document.getElementById('mods').style.display = 'none';
    SetModStatus("Will show boot animation on screen for 10 seconds...")
    let response = await fetch('/api/mods/custom/TestBootAnim', {
        method: 'POST',
    });
    let data = await response.json();
    if (data.status == "success") {
        document.getElementById('mods').style.display = 'block';
        SetModStatus("")
    } else {
        document.getElementById('mods').style.display = 'block';
        SetModStatus("TestBootAnim error: " + data.message)
    }
    return data;
}

function bootAnimCheckValidate() {
    let checkbox = document.getElementById('bootAnimDefault');
    let divUpload = document.getElementById('bootAnimUploadHide');

    if(checkbox.checked == true) {
        divUpload.style.display = "none";
        document.getElementById('bootAnimCurrent').style.display = "none";
    } else {
        divUpload.style.display = "block";
        document.getElementById('bootAnimCurrent').style.display = "block";
    }
}

async function BootAnim_Submit() {
    let checkbox = document.getElementById('bootAnimDefault');
    let inputFile = document.getElementById('bootAnimUpload');
    let gifData = "";

    if(checkbox.checked == false && inputFile.files.length > 0) {
        let file = inputFile.files[0];
        gifData = await new Promise((resolve) => {
            let reader = new FileReader();
            reader.onload = (event) => resolve(event.target.result.split(',')[1]);
            reader.readAsDataURL(file);
        });
    }

    let json = `{"default": ${checkbox.checked}, "gifdata": "${gifData}"}`;
    let banimresp = await SendJSON('BootAnim', json);
    if (banimresp == "error") {
        alert(banimresp)
    }
}







window.onload = async function() {
    let data = await GetCurrent('FreqChange');
    let radioButtons = document.getElementsByName("frequency");
    for(let i = 0; i < radioButtons.length; i++){
        if(radioButtons[i].value == data.freq){
            radioButtons[i].checked = true;
            break;
        }
    }
}

async function GetCurrent(mod) {
    let response = await fetch('/api/mods/current/' + mod);
    let data = await response.json();
    return data;
}

async function SendJSON(mod, json) {
    let response = await fetch('/api/mods/modify/' + mod, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: json,
    });
    let data = await response.json();
    return data;
}

async function FreqChange_Submit() {
    let freq = document.querySelector('input[name="frequency"]:checked').value;
    let data = await SendJSON('FreqChange', `{"freq":` + freq + `}`);
    console.log('Success:', data);
}

class StarLamp {
    constructor() {
        this.baseURL = "http://192.168.1.171:8080"
    }

    async getStatus() {
        let response = await fetch(this.baseURL + "/status")
        return await response.json()
    }

    fillStatus(status) {
        document.getElementById("wake_time").value = status.awakeTime;
        document.getElementById("wake_color").value = status.awakeColor;
        document.getElementById("sleep_time").value = status.asleepTime;
        document.getElementById("sleep_color").value = status.asleepColor;
        document.getElementById("current_status").value = status.currentState;
        document.getElementById("lamp_color").value = status.currentColor;
    }

    fillColorSelects() {
        let selects = document.getElementsByClassName("lamp_color");
        for (let i = 0; i < selects.length; i++) {
            let select = selects[i];
            for (let colorIdx in LightState) {
                if (LightState.hasOwnProperty(colorIdx)) {
                    let color = LightState[colorIdx];
                    select.options.add(new Option(color.name, color.id))
                }
            }
        }
    }

    fillStateSelects() {
        let selects = document.getElementsByClassName("awake_state");
        for (let i = 0; i < selects.length; i++) {
            let select = selects[i];
            for (let stateIdx in AwakeState) {
                if (AwakeState.hasOwnProperty(stateIdx)) {
                    let state = AwakeState[stateIdx];
                    select.options.add(new Option(state.name, state.id))
                }
            }
        }
    }
}

const LightState = {
    "none": {id: 0, name: "All Off"},
    "blue": {id: 1, name: "Blue"},
    "red": {id: 2, name: "Red"},
    "green": {id: 3, name: "Green"},
    "red_blue": {id: 4, name: "Red and Blue"},
    "red_green": {id: 5, name: "Red and Green"},
    "blue_green": {id: 6, name: "Blue and Green"},
    "all": {id: 7, name: "All On"},
    "strobe": {id: 8, name: "Strobe"},
};

const AwakeState = {
    "unknown": {id: 0, name: "Unknown"},
    "awake": {id: 1, name: "Awake"},
    "asleep": {id: 2, name: "Asleep"},
};

document.addEventListener("DOMContentLoaded", async () => {
    let lamp = new StarLamp();
    lamp.fillColorSelects();
    lamp.fillStateSelects();
    let status = await lamp.getStatus();
    lamp.fillStatus(status);
    console.log(status);
});
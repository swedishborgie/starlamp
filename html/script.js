class StarLamp {
    constructor() {
        this.baseURL = "http://192.168.1.171:8080";
    }

    async getState() {
        let response = await fetch(this.baseURL + "/status");
        let status = new LampState(await response.json());
        this.fillStatus(status);
        return status
    }

    async setState(state) {
        state = state || this.buildState();
        let response = await fetch(this.baseURL + "/status", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(state),
        });
        let status = new LampState(await response.json());
        this.fillStatus(status);
        return status;
    }

    async reset() {
        await fetch(this.baseURL + "/reset", {method: "POST"})
        await this.getState()
    }

    async nextColor() {
        let state = await this.getState();
        state.currentColor = (state.currentColor + 1) % 9;
        await this.setState(state)
    }

    fillStatus(status) {
        document.getElementById("wake_time").value = status.awakeTime;
        document.getElementById("wake_color").value = status.awakeColor;
        document.getElementById("sleep_time").value = status.asleepTime;
        document.getElementById("sleep_color").value = status.asleepColor;
        document.getElementById("current_status").value = status.currentState;
        document.getElementById("lamp_color").value = status.currentColor;
    }

    buildState() {
        let state = new LampState();
        state.awakeTime = document.getElementById("wake_time").value;
        state.awakeColor = parseInt(document.getElementById("wake_color").value);
        state.asleepTime = document.getElementById("sleep_time").value;
        state.asleepColor = parseInt(document.getElementById("sleep_color").value);
        state.currentColor = parseInt(document.getElementById("lamp_color").value);
        return state;
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

class LampState {
    constructor(obj) {
        obj = obj || {};
        this.awakeTime = obj.awakeTime || "00:00:00";
        this.awakeColor = obj.awakeColor || 0;
        this.asleepTime = obj.asleepTime || "00:00:00";
        this.asleepColor = obj.asleepColor || 0;
        this.currentState = obj.currentState || 0;
        this.currentColor = obj.currentColor || 0;
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
    await lamp.getState();

    document.getElementById("save").addEventListener("click", () => lamp.setState());
    document.getElementById("next_color").addEventListener("click", () => lamp.nextColor());
    document.getElementById("reset").addEventListener("click", () => lamp.reset());
});
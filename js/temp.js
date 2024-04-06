import { WebSocketController, ViewMode } from './common.js'

export function run(prefix, url, viewMode) {
	const temp = new Temp(prefix, url, viewMode)
}

class Temp extends WebSocketController {

	open() {
		super.open()
		if (this.state.DeployParams !== "") {
			this.showTemp()
		}
	}

	handle(msg) {
		switch(msg.Path) {
		case "update":
			this.update(msg)
			break
		}
	}

	showTemp() {
		let temp = document.getElementById("temp")
		let hum = document.getElementById("hum")

		temp.textContent = this.state.Dht.Temperature.toFixed(1)
		hum.textContent = this.state.Dht.Humidity.toFixed(1)
	}

	update(msg) {
		let dht = this.state.Dht
		dht.Temperature = msg.Temperature
		dht.Humidity = msg.Humidity
		this.showTemp()
	}
}

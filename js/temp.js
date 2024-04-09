import { WebSocketController, ViewMode } from './common.js'

export function run(prefix, url, viewMode) {
	const temp = new Temp(prefix, url, viewMode)
}

class Temp extends WebSocketController {

	open() {
		super.open()
		if (this.state.DeployParams !== "") {
			this.showTemp()
			this.showChart()
		}
	}

	handle(msg) {
		switch(msg.Path) {
		case "update":
			this.update(msg)
			break
		}
	}

	initChart() {
		let canvas = document.getElementById("chartCanvas")
		if (typeof this.chart !== 'undefined') {
			return
		}
		this.chart = new Chart(canvas, {
			type: 'line',
			data: {
				labels: [],
				datasets: [{
					label: 'Temperature',
					data: [],
					yAxisID: 'y',
					borderColor: 'red',
					backgroundColor: 'red',
				}, {
					label: 'Humidity',
					data: [],
					yAxisID: 'y1',
					borderColor: 'blue',
					backgroundColor: 'blue',
				}],
			},
			options: {
				animation: false,
				maintainAspectRatio: false,
				plugins: {
					legend: {
					    display: false,
					},
				},
				scales: {
				    x: {
					title: {
						display: true,
						text: "Last Hour",
					},
				    },
				    y: {
					type: 'linear',
					position: 'left',
					suggestedMin: 0,
					suggestedMax: 30,
					title: {
					    display: true,
					    text: "Temperature °C",
					    color: "red",
					},
				    },
				    y1: {
					type: 'linear',
					position: 'right',
					min: 0,
					max: 100,
					title: {
					    display: true,
					    text: "Humidity %",
					    color: "blue",
					},
					grid: {
          				    drawOnChartArea: false, // only want the grid lines for one axis to show up
        				},

				    },
				},
			},
		})
	}

	showChart() {

		if (this.viewMode !== ViewMode.ViewFull) {
			return
		}

		this.initChart()

		this.chart.data.labels = Array(60).fill("")
		this.chart.data.datasets[0].data = Array(60).fill(null)
		this.chart.data.datasets[1].data = Array(60).fill(null)

		for (let i = 0; i < this.state.History.length; i++) {
			let rec = this.state.History[i]
			this.chart.data.datasets[0].data[60 - i] = rec[0]
			this.chart.data.datasets[1].data[60 - i] = rec[1]
		}

		this.chart.update()
	}

	showTemp() {
		let temp = document.getElementById("temp")
		let hum = document.getElementById("hum")

		temp.textContent = this.state.Dht.Temperature.toFixed(1)
		hum.textContent = this.state.Dht.Humidity.toFixed(1)
	}

	updateHistory() {
		var history = this.state.History
		let dht = this.state.Dht
		if (history.length > 60) {
			history.pop()
		}
		history.unshift([dht.Temperature, dht.Humidity])
	}

	update(msg) {
		let dht = this.state.Dht
		dht.Temperature = msg.Temperature
		dht.Humidity = msg.Humidity
		this.updateHistory()
		this.showTemp()
		this.showChart()
	}
}

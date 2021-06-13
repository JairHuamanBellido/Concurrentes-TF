let monoxidoCarbono =  document.getElementById("monoxido-carbono")
let acidoSulfhridico =  document.getElementById("acido-sulfhidrico")
let dioxidoNitrogeno  = document.getElementById("dioxido-nitrogeno")
let ozono = document.getElementById("ozono")
let pm10 = document.getElementById("pm10")
let pm25 = document.getElementById("pm25")
let dioxidoAzufre = document.getElementById("dioxido-azufre")
let ruido = document.getElementById("ruido")
let uv = document.getElementById("uv")
let humedad = document.getElementById("humedad")
let presion  = document.getElementById("presion")


let submitButton = document.getElementById("submit-btn")
let prediction_text = document.getElementById("prediction-result")


submitButton.addEventListener("click", () => {

    fetch(`http://localhost:3030/?monoxidoCarbono=${monoxidoCarbono.value}&acidoSulfridico=${acidoSulfhridico.value}&dioxidoDeNitrogeno=${dioxidoNitrogeno.value}&ozono=${ozono.value}&pm10=${pm10.value}&pm25=${pm25.value}&dioxidoDeAzufre=${dioxidoAzufre.value}&ruido=${ruido.value}&uv=${uv.value}&humedad=${humedad.value}&presion=${presion.value}`)
    .then((s) => s.json())
    .then((data) => {
        prediction_text.textContent= data.Temperatura
    })
}, false)


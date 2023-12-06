<template>
    <div class="container mt-5">
        <h1 class="mb-4">Weather App v0.0.1 - bazz - 2</h1>
        <div class="mb-3">
            <label for="location" class="form-label">Enter Location:</label>
            <input type="text" id="location" v-model="location" class="form-control" @keyup.enter="getWeather">
            <button @click="getWeather" :disable="loading" class="btn btn-primary mt-2">Get Weather</button>
        </div>
        <div v-if="loading" class="alert alert-info">Loading...</div>
        <div v-if="error" class="alert alert-danger">{{ error }}</div>
        <div v-if="weatherData" class="mt-4">
            <h2>{{ weatherData.city }}, {{ weatherData.country }}</h2>
            <p class="mb-2">Temperature: {{ weatherData.temperature }}°C</p>
            <p class="mb-2">Weather: {{ weatherData.description }}</p>
            <p class="mb-1">web-api v{{ weatherData.webApiVersion }}</p>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
export default {
    data() {
        return {
            location: '',
            weatherData: null,
            loading: false,
            error: null
        };
    },
    methods: {
        getWeather() {
            if (this.location.trim() === '') {
                this.error = 'Please enter a location';
                return;
            }

            this.loading = true;
            this.error = null;

            const apiUrl = `./api/weather?location=${encodeURIComponent(this.location)}`;

            axios.get(apiUrl)
                .then(response => {
                    this.weatherData = response.data;
                })
                .catch(error => {
                    console.error('Error fetching weather data', error);
                    this.error = 'Error fetching weather data';
                })
                .finally(() => {
                    this.loading = false;
                });
        },
    },
};
</script>

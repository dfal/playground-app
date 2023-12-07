import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"

import { createApp } from 'vue'
import WeatherApp from './WeatherApp.vue'

var app = createApp(WeatherApp);
app.config.globalProperties.version = import.meta.env.PACKAGE_VERSION;
app.mount('#app');

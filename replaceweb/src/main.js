import { createApp } from 'vue';
import VueChartkick from 'vue-chartkick';
import 'chartkick/chart.js';
import Autocomplete from '@trevoreyre/autocomplete-vue';
import Moment from 'vue-moment';
import '@trevoreyre/autocomplete-vue/dist/style.css';
import 'flag-icons/css/flag-icons.min.css';

import App from './App.vue'
import router from './router'
import './assets/main.css'

const app = createApp(App)

app.use(router)
app.use(VueChartkick);
app.use(Autocomplete);
// app.use(Moment);

app.mount('#app')

import Vue from 'vue';
import VueRouter from 'vue-router';
import Chartkick from 'vue-chartkick';
import Chart from 'chart.js';
import Autocomplete from '@trevoreyre/autocomplete-vue';
import '@trevoreyre/autocomplete-vue/dist/style.css';
import 'flag-icon-css/css/flag-icon.css';

import App from './App.vue';
import Homepage from './pages/Homepage';
import About from './pages/About';
import './global.css';

Vue.use(Chartkick.use(Chart));
Vue.use(VueRouter);
Vue.use(Autocomplete);

const router = new VueRouter({
    routes: [
        { path: '/', component: Homepage },
        { path: '/about', component: About }
    ]
});

Vue.config.productionTip = false;

new Vue({
    router,
    render: h => h(App)
}).$mount('#app');
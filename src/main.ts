import './assets/main.css';

import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { createVuetify } from 'vuetify';
import 'vuetify/styles';

const vuetify = createVuetify({
  defaults: {
    global: {
      variant: 'outlined'
    }
  }
});

const app = createApp(App);
app.use(vuetify);
app.use(router);
app.mount('#app');

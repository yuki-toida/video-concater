import Vue from 'vue';
import app from './app.vue';
import router from './router';

new Vue({
  el: '#app',
  router,
  components: { app },
  render: h => h(app),
});

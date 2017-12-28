// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
//import router from './router'

Vue.config.productionTip = false

/* eslint-disable no-new */
var app = new Vue({
  el: '#app',
//  router,
//  template: '<App/>',
  data: {
    message: "hello world"
  },
//  components: { App }
})

var app2 = new Vue({
  el: "#app-2",
  data: {
    message: "fuck you"
  }
})

var app4 = new Vue({
  el: "#app-4",
  data: {
    todos: [
      {text: "haha"},
      {text: "haha1kj"},
      {text: "haha1kj2222"}
    ]
  }
})

var app5 = new Vue({
  el: "#app5",

  data: {
    message: "em.....",
    count: 0,
    seen: false,
    disabled: true
  },

  computed: {
    now: function ()  {
      return Date.now()
    }
  },

  methods: {
    onMe: function () {
      this.message = 'why click me!!!!!!' + this.count
      this.count += 1
    }
  }

})

Vue.compoment('todo-item', {
  template: '<li>This is todo</li>'
})



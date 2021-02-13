<template>
  <div class="hello">
    <!-- <h1>{{ msg }}</h1> -->
    
    <div>
      <h1 class="title">Configuration Values</h1>
    </div>
    <hr>
    
          <div class="form-group">
            <div class="field">
              <label class="label">Port Number</label>
              <input v-model="portNumber" type="text" placeholder="Enter the port number" class="form-control">
            </div>
            <div class="field" style="margin-top:13px;">
              <label class="label">Message</label>
              <input v-model="message" type="text"  placeholder="Enter the message" class="form-control">
            </div>
          </div>
          <div class="form-group">
            <button class="btn btn-primary" v-on:click="applyChanges()">Apply</button>
          </div>

    &nbsp;
    <hr>

  </div>
</template>

<script>

import axios from 'axios';


/* eslint-disable */


export default {
  name: 'Controller',
  
  data: function() {
    return {
    portNumber: '',
    message: '',
  }
  },
  created () {
    // fetch the data when the view is created and the data is
    // already being observed
    this.getInitialValues()
  },
  methods: {
    applyChanges: function() {
      var data = {"portNumber": parseInt(this.portNumber), "message": this.message}

      axios({ method: "POST", url: "http://localhost:8090/", data: data, headers: {"content-type": "text/plain" } }).then(result => { 

        }).catch( error => {
            /*eslint-disable*/
            window.alert(`${error} ${this.portNumber} ${this.message} `);
            /*eslint-enable*/
      });
    },
    getInitialValues: function() {
      axios({ method: "GET", url: "http://localhost:8090/", headers: {"content-type": "text/plain" } }).then(result => { 
            this.portNumber = result.data.portNumber;
            this.message = result.data.message;

        }).catch( error => {
            /*eslint-disable*/
            window.alert(`${error} ${this.portNumber} ${this.message} `);
            /*eslint-enable*/
      });
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>

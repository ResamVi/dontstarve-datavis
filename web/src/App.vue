<template>
  <div id="app">
    <nav>
      <img src="./assets/logo.png" id="logo">
      <h1 id="title">ResamVi's DST DataViz</h1>
      <ul id="navigation">
        <li><router-link tag="a" to="/">Homepage</router-link></li>
        <li><router-link tag="a" to="/about">About</router-link></li>
      </ul>
    </nav>
    
    <router-view></router-view>    
  </div>
</template>

<script>

export default {
  name: 'App',
  components: {
    
  }
}
</script>

<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #5b4421;
}

#title {
  vertical-align: middle;
  display: inline-block;
  margin-right: 10px;
  color: #f7f2ea;
  margin: 0;
  font-weight: 100;
}

#logo {
  vertical-align: middle;
  display: inline-block;
  margin-right: 10px;
  height: 34px;
  width: 34px;
}

#navigation {
  display: inline-block;
  vertical-align: middle;
  list-style-type: none;
  color: #f7f2ea;
  margin: 0;
  padding: 0;
}

#navigation li {
  float: left;
  margin: 10px;
}

#navigation li a {
  display: block;
  padding: 10px 30px;
  color: white;
  text-align: center;
  text-decoration: none;
}

#navigation li a:hover {
  background-color: #111;
}

</style>

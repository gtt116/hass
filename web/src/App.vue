<template>
  <div id="app">

<el-container>
  <el-header>
  <a href="https://github.com/gtt116/hass">
    <img src="/static/favicon.png" width="40"/>
  </a>
  </el-header>

  <el-container>
    <el-aside width="200px">
      <el-menu default-active="/" :router="true">
        <el-menu-item index="/">
          <i class="el-icon-menu"></i>
          <span slot="title">Dashboard</span>
        </el-menu-item>
        <el-menu-item index="/setting">
          <i class="el-icon-setting"></i>
          <span slot="title">Setting</span>
        </el-menu-item>
      </el-menu> 
    </el-aside>

    <el-container>
      <el-main>
        <router-view/>
      </el-main>
    </el-container>

  </el-container>
</el-container>

  </div>
</template>

<script>
export default {
    name: 'app'
}
</script>

<style>
  body {
    margin: 0px;
    font-family: "Helvetica Neue",Helvetica,"PingFang SC","Hiragino Sans GB","Microsoft YaHei","微软雅黑",Arial,sans-serif;
  }
  a {
      text-decoration: none;
  }
  .el-header {
    background-color: #4267b2;
    text-align: left;
    line-height: 88px;
  }
  
  .el-aside {
    background-color: #e6e6e6;
    color: #333;
    text-align: left;
    line-height: 200px;
  }
  
  .el-main {
    min-height: 200px;
  }

.detail {
  border: 1px solid #eee;
  padding:10px 30px 30px;
  margin-bottom: 40px;
}
</style>

<template>
<div>
  <h2 class="title">
    <span>Dashboard</span>
    <el-button type="primary" size="mini" class="title-button" @click="refresh">
      <i class="el-icon-refresh"></i>
    </el-button> 
  </h2>
  <div id="wall">
  <el-row :gutter="40">
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-upload"></i>
        <span>{{ count }}</span>
        Servers
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-rank"></i>
        <span>{{ connections }} </span>
        Connections
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-upload2"></i>
        <span>{{ sent }} </span>
        KB/s
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-download"></i>
        <span>{{ recv }}</span>
        KB/s
      </div>
    </el-col>
  </el-row>
  </div>

  <el-row class="detail">
    <h4>Servers</h4>
    <div>
      <el-table :data="servers" stripe>
        <el-table-column prop="ip" label="IP">
        </el-table-column>
        <el-table-column prop="sent" label="Sent Rate">
        </el-table-column>
        <el-table-column prop="recv" label="Recv Rate">
        </el-table-column>
        <el-table-column prop="connections" label="Connections">
        </el-table-column>
        <el-table-column prop="trend" label="Trends">
        </el-table-column>
      </el-table>
    </div>
  </el-row>

  <el-row class="detail">
    <h4>Connections</h4>
    <div>
      <el-table :data="connection_list" stripe>
        <el-table-column prop="source" label="source"></el-table-column>
        <el-table-column prop="server" label="server"></el-table-column>
        <el-table-column prop="target" label="target"></el-table-column>
        <el-table-column prop="sent" label="sent"> </el-table-column>
        <el-table-column prop="recv" label="recv"> </el-table-column>
      </el-table>
    </div>
  </el-row>

</div>
</template>

<script>
export default {
  name: 'Dashboard',
  methods: {
    refresh () {
      console.log("refresh")
    }
  },
  data () {
    return {
      count: 5,
      connections: 5,
      sent: 30,
      recv: 39,
      servers: [{
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        connections:123,
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        connections:123,
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        connections:123,
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        connections:123,
        trend: [12,34,35,54]
      }],
      connection_list: [{
        source: "127.0.0.1:232",
        server: "2.3.4.5:242",
        target: "www.google.com:443",
        sent: 234,
        recv: 234,
      }, {
        source: "127.0.0.1:232",
        server: "2.3.4.5:242",
        target: "www.google.com:443",
        sent: 234,
        recv: 234,
      }, {
        source: "127.0.0.1:232",
        server: "2.3.4.5:242",
        target: "www.google.com:443",
        sent: 234,
        recv: 234,
      }
      ]
    }
  }
}
</script>

<style>
#wall {
  margin-bottom: 40px;
}
.hero {
  text-align: center;
  background-color: #ffffff;
  line-height: 100px;
  font-size: 20px;
  border: 1px solid #eee;
  border-radius: 4px;
  box-shadow: 1px 1px 1px #eee;
}

.hero span {
  font-weight: bold;
  font-size: 26px;
}
</style>


<template>
<div>
  <div id="wall">
  <el-row :gutter="40">
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-upload"></i>
        {{ total.count }} Servers
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-rank"></i>
        {{ total.connections }} Connections
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-upload2"></i>
        {{ total.sent }} KB/s
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-download"></i>
        {{ total.recv }} KB/s
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
        <el-table-column prop="sent" label="Sent rate">
        </el-table-column>
        <el-table-column prop="recv" label="Recv rate">
        </el-table-column>
        <el-table-column prop="msg" label="Messge">
        </el-table-column>
        <el-table-column prop="trend" label="trends">
        </el-table-column>
      </el-table>
    </div>
  </el-row>

  <el-row class="detail">
    <h4>Connections</h4>
    <div>
      <el-table :data="connections" stripe>
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
  data () {
    return {
      total: {
        count: 5,
        connections: 5,
        sent: 30,
        recv: 39
      },
      servers: [{
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        msg: "Hi loow",
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        msg: "Hi loow",
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        msg: "Hi loow",
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        msg: "Hi loow",
        trend: [12,34,35,54]
      }],
      connections: [{
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
      }]
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
  background-color: #E9EEF3;
  line-height: 100px;
  font-size: 20px;
}
</style>


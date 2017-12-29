<template>
<div>
  <div id="wall">
  <el-row :gutter="40">
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-upload"></i>
        {{ count }} Servers
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-rank"></i>
        {{ connections }} Connections
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-upload2"></i>
        {{ sent }} KB/s
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-download"></i>
        {{ recv }} KB/s
      </div>
    </el-col>
  </el-row>
  </div>

  <el-row class="detail">
    <h4>Servers</h4>
    <div>
      <el-table :data="tableData" stripe>
        <el-table-column prop="ip" label="IP" width="100">
        </el-table-column>
        <el-table-column prop="sent" label="sent rate" width="100">
        </el-table-column>
        <el-table-column prop="recv" label="recv rate" width="200"> 
        </el-table-column>
        <el-table-column prop="trend" label="trends">
        </el-table-column>
      </el-table>
    </div>
  </el-row>

  <el-row class="detail">
    <h4>Connections</h4>
    <div>
      <el-table :data="tableData2" stripe>
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
      count: 5,
      connections: 5,
      sent: 30,
      recv: 39,
      tableData: [{
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        trend: [12,34,35,54]
      }, {
        ip: "2.3.4.5",
        sent: 23,
        recv: 23,
        trend: [12,34,35,54]
      }],
      tableData2: [{
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
  background-color: #E9EEF3;
  line-height: 100px;
  font-size: 20px;
}
</style>


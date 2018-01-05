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
        <span>{{ total.count }}</span> Servers
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-rank"></i>
        <span>{{ total.connections }}</span> Connections
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-upload2"></i>
        <span>{{ total.sent }}</span> KB/s
      </div>
    </el-col>
    <el-col :span="6">
      <div class="hero">
        <i class="el-icon-download"></i>
        <span>{{ total.recv }}</span> KB/s
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
  methods: {
    refresh () {
      this.fetchData()
    },
    fetchData () {
      fetch('http://127.0.0.1:7777/api/total')
      .then(res => { return res.json() })
      .then(json => {
        this.total = json
      })

      fetch('http://127.0.0.1:7777/api/servers')
      .then(res => { return res.json() })
      .then(json => {
        this.servers = json
      })

      fetch('http://127.0.0.1:7777/api/connections')
      .then(res => { return res.json() })
      .then(json => {
        this.connections = json
      })
    }
  },
  created () {
    this.fetchData()
  },
  watch: {
    '$route': 'fetchData'
  },
  data () {
    return {
      total: { },
      servers: [],
      connections: []
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


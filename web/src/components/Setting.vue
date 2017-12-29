<template>
<div>
  <el-row class="detail">
    <h3>Basic</h3>
    <div class="setting-form">
      <el-form ref="form" :model="form" label-width="100px">
        <el-form-item label="Host">
          <el-input v-model="form.host"></el-input>
        </el-form-item>
        <el-form-item label="AdminPort">
          <el-input v-model="form.admin_port"></el-input>
        </el-form-item>
        <el-form-item label="Socks5Port">
          <el-input v-model="form.socks_port"></el-input>
        </el-form-item>
        <el-form-item label="HttpPort">
          <el-input v-model="form.http_port"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submit1">Save</el-button>
          <el-button>Cancel</el-button>
        </el-form-item>
      </el-form> 
    </div>
    </el-row>

  <el-row class="detail">
    <div class="setting-form">
      <h3>默认配置</h3>
      <el-form ref="form2" :inline="true" :model="form2" label-width="100px">
        <el-form-item label="Port">
          <el-input v-model="form2.default_port"></el-input>
        </el-form-item>
        <el-form-item label="Password">
          <el-input type="password" auto-complete="off" v-model="form2.default_password"></el-input>
        </el-form-item>
        <el-form-item label="Method">
          <el-select v-model="form2.default_method" placeholder="default method">
            <el-option label="rc4-md5" value="rc4-md5"></el-option>
            <el-option label="shasha" value="shasha"></el-option>
          </el-select>
        </el-form-item>
      </el-form> 

      <h3>Details</h3>
      <el-form ref="form3" :model="form3" label-width="100px">
        <el-form-item 
        v-for="(server, index) in form3.servers"
        :label="index + 1" class='detail-form' :key="server.key" 
        >
        <el-row>
          <el-input v-model="server.ip" placeholder="ip"></el-input>
          <el-input v-model="server.port" clearable placeholder="port"></el-input>
          <el-input v-model="server.password" clearable type="password" placeholder="password"></el-input>
          <el-select v-model="server.method" placeholder="method">
            <el-option label="rc4-md5" value="rc4-md5"></el-option>
            <el-option label="shasha" value="shasha"></el-option>
          </el-select>
          <el-button type="danger" size="medium" @click.prevent="removeIt(server)"><i class="el-icon-delete"></i></el-button>
        </el-row>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submit2">Save</el-button>
          <el-button @click="addMore">Add</el-button>
          <el-button>Cancel</el-button>
        </el-form-item>
      </el-form> 
    </div>

  </el-row>

</div>
</template>

<script>
export default {
  name: 'Setting',
  methods: {
    submit1() {
      console.log(this.form)
    },
    submit2() {
      console.log(this.form2);
      console.log(this.form3);
    },
    addMore () {
        this.form3.servers.push({})
    },
    removeIt(item) {
        this.form3.servers = this.form3.servers.filter(s => s != item)
    }
  },
  data () {
    return {
      form: {
        host: '',
        admin_port: 123,
        socks_port: 123,
        http_port: 23
      },
      form2: {
      },
      form3: {
        servers: [
            {
                ip: '1.2.43'
            }, 
            {
                ip: '2.3.4.5'
            }, 
            {
                ip: '4.3.4.5'
            }
        ]
      }
    }
  }
}
</script>

<style>
.setting-form {
  width: 90%;
}

.el-input {
  width: auto;
}

</style>

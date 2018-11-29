<template>
  <div>
    <!-- <el-col :span="24">

    </el-col> -->
    <div>FAKE API SERVER</div>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-button type="primary"
                   @click="dialogVisible = true">Add API</el-button>
      </el-col>

      <el-col :span="12"></el-col>
    </el-row>
    <br>
    <el-row :gutter="20">

      <el-col :span="3">&nbsp;</el-col>
      <el-col :span="18">

        <el-table :data="tableData"
                  border
                  style="width: 100%">

          <el-table-column prop="id"
                           label="ID"
                           width="180">
          </el-table-column>
          <el-table-column prop="protocol"
                           label="Protocol"
                           width="180">
          </el-table-column>
          <el-table-column prop="url"
                           label="URL"
                           width="Response">
          </el-table-column>
          <el-table-column prop="responseStatus"
                           label="Response Status">
          </el-table-column>
          <el-table-column prop="responseStr"
                           label="Response String">
          </el-table-column>

          <el-table-column label="Operation">
            <template slot-scope="scope">
              <el-button size="mini"
                         type="success"
                         @click="clickEdit(scope.row)">Edit</el-button>
              <el-button size="mini"
                         type="danger"
                         @click="handleDelete(scope.row)">Delete</el-button>
            </template>
          </el-table-column>

        </el-table>

      </el-col>
      <el-col :span="3">&nbsp;</el-col>
    </el-row>
    <el-col :span="24">

      <el-dialog title="Add Fake API"
                 :visible.sync="dialogVisible"
                 width="60%">
        <el-alert title="Empty Content"
                  type="error"
                  show-icon="true"
                  description="You have to input Protocol, URL and Response(status, string)!!"
                  v-show="showAlert">
        </el-alert>

        <span>Protocol</span>
        <el-select v-model="inputProtocol"
                   placeholder="Protocol">
          <el-option v-for="item in protocolOptions"
                     :key="item.value"
                     :label="item.label"
                     :value="item.value">
          </el-option>
        </el-select>

        <br>
        <br>
        <span>Response Status</span>
        <el-select v-model="inputResponseStatus"
                   placeholder="Response Status">
          <el-option v-for="item in statusOptions"
                     :key="item.value"
                     :label="item.label"
                     :value="item.value">
          </el-option>
        </el-select>
        <br> <br>

        <el-input v-model="inputUrl"
                  placeholder="URL, ex: /api/test">
        </el-input>

        <br> <br>

        <el-input type="textarea"
                  v-model="inputResponseStr"
                  rows="8"
                  placeholder="Response, ex: {}"></el-input>

        <span slot="footer"
              class="dialog-footer">
          <el-button @click="dialogVisible = false">Cancel</el-button>
          <el-button type="primary"
                     @click="handleConfirm">Confirm</el-button>
        </span>
      </el-dialog>
    </el-col>

    <el-col :span="24">
      <el-dialog title="Edit Fake API"
                 :visible.sync="editDialogVisible"
                 width="60%">
        <span>Protocol</span>
        <el-select v-model="editProtocol"
                   placeholder="Protocol">
          <el-option v-for="item in options"
                     :key="item.value"
                     :label="item.label"
                     :value="item.value">
          </el-option>
        </el-select>

        <br>
        <br>
        <span>Response Status</span>
        <el-select v-model="editResponseStatus"
                   placeholder="Response Status">
          <el-option v-for="item in statusOptions"
                     :key="item.value"
                     :label="item.label"
                     :value="item.value">
          </el-option>
        </el-select>
        <br> <br>

        <el-input v-model="editUrl"
                  placeholder="URL, ex: /api/test">
        </el-input>

        <br> <br>

        <el-input type="textarea"
                  v-model="editResponseStr"
                  rows="8"
                  placeholder="Response, ex: {}"></el-input>

        <span slot="footer"
              class="dialog-footer">
          <el-button @click="editDialogVisible = false">Cancel</el-button>
          <el-button type="primary"
                     @click="handleEdit()">Confirm</el-button>
        </span>
      </el-dialog>
    </el-col>

  </div>
</template>

<script>
import qs from 'qs'
let domain = ''

function isEmpty (str) {
  return (!str || str.length === 0)
}

export default {
  name: 'list',
  data () {
    return {
      dialogVisible: false,
      editDialogVisible: false,
      tableData: [],
      protocolOptions: [{
        value: 'GET',
        label: 'GET'
      }, {
        value: 'POST',
        label: 'POST'
      }, {
        value: 'PUT',
        label: 'PUT'
      }, {
        value: 'DELETE',
        label: 'DELETE'
      }],
      statusOptions: [{
        value: 200,
        label: 200
      }, {
        value: 400,
        label: 400
      }, {
        value: 401,
        label: 401
      }, {
        value: 404,
        label: 404
      }, {
        value: 409,
        label: 409
      }, {
        value: 500,
        label: 500
      }, {
        value: 502,
        label: 502
      }, {
        value: 503,
        label: 503
      }],
      inputProtocol: '',
      inputResponseStatus: 200,
      inputUrl: '',
      inputResponseStr: '',
      editId: '',
      editProtocol: '',
      editUrl: '',
      editResponseStr: '',
      editResponseStatus: '',
      showAlert: false
    }
  },
  methods: {

    handleConfirm (done) {
      let serverUrl = domain + '/apicontent'
      let url = this.$data.inputUrl
      if (url.charAt(0) !== '/') {
        url = '/' + url
      }

      let formData = {
        protocol: this.$data.inputProtocol,
        url: url,
        responseStr: this.$data.inputResponseStr,
        responseStatus: this.$data.inputResponseStatus
      }
      if (isEmpty(this.$data.inputProtocol) ||
        isEmpty(this.$data.inputUrl) ||
        isEmpty(this.$data.inputResponseStr)) {
        this.$data.showAlert = true
        return
      }
      this.$data.dialogVisible = false
      this.axios.post(serverUrl, qs.stringify(formData))
        .then(res => {
          this.axios.get(domain + '/apicontent').then(res => {
            this.$data.tableData = res.data
          })
        })
    },
    clickEdit (row) {
      this.$data.editId = row.id
      this.$data.editProtocol = row.protocol
      this.$data.editUrl = row.url
      this.$data.editResponseStr = row.responseStr
      this.$data.editResponseStatus = row.responseStatus

      this.$data.editDialogVisible = true
    },
    handleEdit () {
      this.$data.editDialogVisible = false
      let serverUrl = domain + '/apicontent/' + this.$data.editId
      let url = this.$data.editUrl
      if (url.charAt(0) !== '/') {
        url = '/' + url
      }
      let formData = {
        protocol: this.$data.editProtocol,
        url: url,
        responseStr: this.$data.editResponseStr,
        responseStatus: this.$data.editResponseStatus
      }

      this.axios.put(serverUrl, qs.stringify(formData)).then(res => {
        this.axios.get(domain + '/apicontent').then(res => {
          this.$data.tableData = res.data
        })
      })
    },
    handleDelete (row) {
      this.axios.delete(domain + '/apicontent/' + row.id).then(res => {
        this.axios.get(domain + '/apicontent').then(res => {
          this.$data.tableData = res.data
        })
      })
    }
  },
  created () {
    this.axios.get(domain + '/apicontent').then(res => {
      this.$data.tableData = res.data
    })
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
  font-weight: normal;
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

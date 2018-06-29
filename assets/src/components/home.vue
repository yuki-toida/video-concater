<template>
  <div class="container">
    <div v-if="uid">
      <div v-if="isNotification" class="notification">
        <button v-on:click="isNotification = false" class="delete"></button>
        <p>結合が正常に終了し、ダウンロード可能になりました</p>
        <p><strong>保存期間は10分です。</strong></p>
      </div>
      <div class="field">
        <a v-bind:href="link" download="output.mp4" class="button is-success">ダウンロード</a>
      </div>
      <div class="field">
        <button v-on:click="init" class="button is-danger">ファイル削除</button>
      </div>
    </div>
    <div v-else>
      <div v-for="(name, index) in inputs" v-bind:key="index" class="field file has-name">
        <label class="file-label">
          <input class="file-input" type="file" accept="video/*" v-on:change="change($event, index)">
          <span class="file-cta">
            <span class="file-icon">
              <i class="fas fa-upload"></i>
            </span>
            <span class="file-label">
              Choose a file…
            </span>
          </span>
          <span class="file-name">
            {{ name }}
          </span>
        </label>
      </div>
      <div class="field">
        <button
          v-on:click="concat"
          v-bind:class="{ 'is-loading': isLoading }"
          v-bind:disabled="files.length != inputs.length"
          class="button is-primary">ファイル結合</button>
        <span v-if="files.length == inputs.length">
          <button v-on:click="add" class="button is-light">ファイル追加</button>
          <button v-if="1 < files.length" v-on:click="remove" class="button is-warning">ファイル削除</button>
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

const domain = process.env.NODE_ENV === 'development'
  ? 'http://localhost:8080'
  : 'http://concat.theliveup.tv';
let staticUrl = null;

export default {
  data: function() {
    return {
      uid: null,
      inputs: [''],
      files: [],
      isLoading: false,
      isNotification: false,
    }
  },
  created: function() {
    axios.get(`${domain}/init`)
      .then((res) => {
        if (res.data) {
          this.uid = res.data.cookie;
          staticUrl = res.data.staticUrl;
        }
      });
  },
  computed: {
    link: function() {
      return `${staticUrl}/static/outputs/${this.uid}.mp4`;
    },
  },
  methods: {
    init: function() {
      axios.delete(`${domain}/cookie`)
        .then((res) => {
          document.cookie = `${res.data}=; max-age=0`;
          this.uid = null;
        });
    },
    change: function(event, index) {
      const file = event.target.files[0];
      this.files.splice(index, 1, {key: index, value: file});
      this.inputs[index] = file.name;
    },
    add: function() {
      this.inputs.push('');
    },
    remove: function() {
      this.files.splice(-1,1);
      this.inputs.splice(-1,1);
    },
    concat: function() {
      this.isLoading = true;
      let formData = new FormData();
      this.files.forEach(element => {
        formData.append(element.key, element.value);
      });

      axios.post(`${domain}/concat`, formData, {headers: {'Content-Type': 'multipart/form-data'}})
        .then(res => {
          console.log(res.data);
          this.uid = res.data;
          this.inputs = [''];
          this.files = [];
          this.isLoading = false;
          this.isNotification = true;
        })
        .catch(error => alert(error));
    }
  }
}
</script>

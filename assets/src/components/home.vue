<template>
  <div class="container">
    <div v-if="uid">
      <div>
        <a v-bind:href="link" download="output.mp4">ダウンロード</a>
      </div>
      <div>
        <button v-on:click="init">削除</button>
      </div>
    </div>
    <div v-else>
      <div v-for="key in inputs" v-bind:key="key" class="file has-name">
        <label class="file-label">
          <input class="file-input" type="file" accept="video/*" v-on:change="change($event, key)">
          <span class="file-cta">
            <span class="file-icon">
              <i class="fas fa-upload"></i>
            </span>
            <span class="file-label">
              Choose a file…
            </span>
          </span>
          <span class="file-name">
            ほげほげ
          </span>
        </label>
      </div>
      <div v-if="files.length == inputs">
        <button v-on:click="add">ファイル追加</button>
        <button v-if="1 < files.length" v-on:click="remove">ファイル削除</button>
      </div>
      <div>
        <button v-on:click="concat">ファイル結合</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

const domain = 'http://localhost:8080';

export default {
  data: function() {
    return {
      uid: null,
      inputs: 1,
      files: [],
    }
  },
  created: function() {
    axios.get(`${domain}/cookies`)
      .then((res) => {
        if (res.data) {
          this.uid = res.data;
        }
      });
  },
  computed: {
    link: function() {
      return `${domain}/outputs/${this.uid}.mp4`;
    }
  },
  methods: {
    init: function() {
      axios.delete(`${domain}/cookies`)
        .then((res) => {
          document.cookie = `${res.data}=; max-age=0`;
          this.uid = null;
        });
    },
    change: function(event, key) {
      const file = event.target.files[0];
      this.files.splice(key - 1, 1, {key: key, value: file});
    },
    add: function() {
      this.inputs++;
    },
    remove: function() {
      this.files.splice(-1,1);
      this.inputs--;
    },
    concat: function() {
      let formData = new FormData();
      this.files.forEach(element => {
        formData.append(element.key, element.value);
      });

      axios.post(`${domain}/concat`, formData, {headers: {'Content-Type': 'multipart/form-data'}})
        .then(res => {
          this.uid = res.data;
          this.inputs = 1;
          this.files = [];
        })
        .catch(error => alert(error));
    }
  }
}
</script>

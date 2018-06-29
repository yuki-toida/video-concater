<template>
  <div>
    <section class="section">
      <div class="container">
        <h1 class="title">
          <router-link v-bind:to="{ name: 'Home' }">MP4結合ツール!</router-link>
        </h1>
        <h2 class="subtitle">
          MP4ファイルを結合しダウンロードするツール
        </h2>
      </div>
    </section>
    <router-view></router-view>
  </div>
</template>

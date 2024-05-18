<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <v-container>
    <main><RouterView /></main>
    <footer>
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/usage">Usage</RouterLink>
      </nav>
    </footer>
  </v-container>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

v-container {
  vertical-align: top;
}

main {
  min-height: 60vh;
}

footer {
  width: 50vw;
}

nav {
  width: 100%;
  font-size: 12px;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>

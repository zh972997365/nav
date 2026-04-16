<template>
  <section class="admin-page-shell">
    <div class="admin-page-shell__content">
      <div v-if="$slots.toolbar" class="admin-page-shell__toolbar">
        <slot name="toolbar" />
      </div>
      <slot />
    </div>
  </section>
</template>

<script>
export default {
  name: 'AdminPageShell',
  props: {
    eyebrow: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
    description: {
      type: String,
      default: '',
    },
  },
};
</script>

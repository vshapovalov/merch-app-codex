<template>
  <section class="page">
    <header class="flex align-items-center justify-content-between mb-4">
      <h2 class="m-0">Пользователи</h2>
      <Button label="Добавить" icon="pi pi-plus" @click="openCreate" />
    </header>
    <DataTable :value="items" dataKey="id" :loading="loading" responsiveLayout="scroll">
      <Column field="name" header="Имя" sortable />
      <Column field="email" header="Email" sortable />
      <Column header="Действия" style="width: 12rem">
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button icon="pi pi-pencil" severity="info" text rounded @click="openEdit(data)" />
            <Button icon="pi pi-trash" severity="danger" text rounded @click="deleteItem(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" modal :header="dialogTitle" class="w-full md:w-6">
      <form v-if="currentItem" class="flex flex-column gap-3" @submit.prevent="saveItem">
        <span class="p-float-label">
          <InputText id="name" v-model="currentItem.name" required />
          <label for="name">Имя</label>
        </span>
        <span class="p-float-label">
          <InputText id="email" v-model="currentItem.email" type="email" required />
          <label for="email">Email</label>
        </span>
        <span class="p-float-label">
          <Password id="password" v-model="currentItem.password" :feedback="false" toggleMask :required="!currentItem.id" />
          <label for="password">Пароль</label>
        </span>
        <div class="flex justify-content-end gap-2">
          <Button label="Отмена" severity="secondary" text @click="dialogVisible = false" type="button" />
          <Button label="Сохранить" type="submit" :loading="saving" />
        </div>
      </form>
    </Dialog>
  </section>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useCrud } from '../composables/useCrud';

const {
  items,
  loading,
  saving,
  dialogVisible,
  currentItem,
  loadItems,
  openCreate,
  openEdit,
  saveItem: baseSave,
  deleteItem,
} = useCrud('/users', () => ({ name: '', email: '', password: '' }), {
  preparePayload: (payload) => {
    if (!payload.password) {
      delete payload.password;
    }
    return payload;
  },
});

const dialogTitle = computed(() =>
  currentItem?.value?.id ? 'Редактирование пользователя' : 'Новый пользователь'
);

const saveItem = () => baseSave();

onMounted(loadItems);
</script>

<style scoped>
.page {
  background: var(--surface-card);
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: var(--card-shadow, 0 1px 4px rgba(0, 0, 0, 0.1));
}
</style>

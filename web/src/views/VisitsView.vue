<template>
  <section class="page">
    <header class="flex align-items-center justify-content-between mb-4">
      <h2 class="m-0">Визиты</h2>
      <Button label="Добавить" icon="pi pi-plus" @click="openCreate" />
    </header>
    <DataTable :value="items" dataKey="id" :loading="loading" responsiveLayout="scroll">
      <Column header="Пользователь">
        <template #body="{ data }">
          {{ userLabelById[data.user_id] ?? data.user_id ?? '' }}
        </template>
      </Column>
      <Column header="Торговая точка">
        <template #body="{ data }">
          {{ retailLabelById[data.retail_point_id] ?? data.retail_point_id ?? '' }}
        </template>
      </Column>
      <Column field="visited_at" header="Дата посещения">
        <template #body="{ data }">
          {{ formatDate(data.visited_at) }}
        </template>
      </Column>
      <Column field="notes" header="Заметки" />
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
        <div>
          <label class="block mb-2">Пользователь</label>
          <Dropdown
            v-model="currentItem.user_id"
            :options="userOptions"
            optionLabel="name"
            optionValue="id"
            placeholder="Выберите пользователя"
            class="w-full"
            required
          />
        </div>
        <div>
          <label class="block mb-2">Торговая точка</label>
          <Dropdown
            v-model="currentItem.retail_point_id"
            :options="retailOptions"
            optionLabel="name"
            optionValue="id"
            placeholder="Выберите точку"
            class="w-full"
            required
          />
        </div>
        <div>
          <label class="block mb-2">Дата посещения</label>
          <Calendar v-model="currentItem.visited_at" showTime hourFormat="24" class="w-full" required />
        </div>
        <span class="p-float-label">
          <Textarea id="notes" v-model="currentItem.notes" autoResize rows="3" />
          <label for="notes">Заметки</label>
        </span>
        <div class="flex justify-content-end gap-2">
          <Button label="Отмена" type="button" severity="secondary" text @click="dialogVisible = false" />
          <Button label="Сохранить" type="submit" :loading="saving" />
        </div>
      </form>
    </Dialog>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useCrud } from '../composables/useCrud';
import api from '../services/api';

const userOptions = ref([]);
const retailOptions = ref([]);

const userLabelById = computed(() =>
  Object.fromEntries(userOptions.value.map((user) => [user.id, user.name]))
);

const retailLabelById = computed(() =>
  Object.fromEntries(retailOptions.value.map((point) => [point.id, point.name]))
);

const crud = useCrud(
  '/visits',
  () => ({ user_id: null, retail_point_id: null, visited_at: new Date(), notes: '' }),
  {
    preparePayload: (payload) => ({
      ...payload,
      visited_at: payload.visited_at ? new Date(payload.visited_at).toISOString() : null,
    }),
  }
);

const { items, loading, saving, dialogVisible, currentItem, deleteItem, loadItems } = crud;

const dialogTitle = computed(() =>
  currentItem?.value?.id ? 'Редактирование визита' : 'Новый визит'
);

const formatDate = (value) => {
  if (!value) return '';
  const date = typeof value === 'string' ? new Date(value) : value;
  return new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date);
};

const openCreate = () => {
  crud.openCreate();
  if (crud.currentItem.value) {
    crud.currentItem.value.visited_at = new Date();
  }
};

const openEdit = (item) => {
  crud.currentItem.value = {
    ...item,
    visited_at: item.visited_at ? new Date(item.visited_at) : new Date(),
  };
  crud.dialogVisible.value = true;
};

const saveItem = () => crud.saveItem();

const loadOptions = async () => {
  const [usersResponse, retailResponse] = await Promise.all([
    api.get('/users'),
    api.get('/retail-points'),
  ]);
  userOptions.value = usersResponse.data;
  retailOptions.value = retailResponse.data;
};

onMounted(async () => {
  await Promise.all([loadItems(), loadOptions()]);
});
</script>

<style scoped>
.page {
  background: var(--surface-card);
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: var(--card-shadow, 0 1px 4px rgba(0, 0, 0, 0.1));
}
</style>

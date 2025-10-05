<template>
  <section class="page">
    <header class="flex align-items-center justify-content-between mb-4">
      <h2 class="m-0">Продукты</h2>
      <Button label="Добавить" icon="pi pi-plus" @click="openCreate" />
    </header>
    <DataTable :value="items" dataKey="id" :loading="loading" responsiveLayout="scroll">
      <Column field="name" header="Название" sortable />
      <Column field="sku" header="SKU" sortable />
      <Column header="Бренд">
        <template #body="{ data }">
          {{ findLabel(brandOptions, data.brand_id) }}
        </template>
      </Column>
      <Column header="Категория">
        <template #body="{ data }">
          {{ findLabel(categoryOptions, data.category_id) }}
        </template>
      </Column>
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
          <label for="name">Название</label>
        </span>
        <span class="p-float-label">
          <InputText id="sku" v-model="currentItem.sku" required />
          <label for="sku">SKU</label>
        </span>
        <div>
          <label class="block mb-2">Бренд</label>
          <Dropdown
            v-model="currentItem.brand_id"
            :options="brandOptions"
            optionLabel="name"
            optionValue="id"
            placeholder="Выберите бренд"
            class="w-full"
            required
          />
        </div>
        <div>
          <label class="block mb-2">Категория</label>
          <Dropdown
            v-model="currentItem.category_id"
            :options="categoryOptions"
            optionLabel="name"
            optionValue="id"
            placeholder="Выберите категорию"
            class="w-full"
            required
          />
        </div>
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

const brandOptions = ref([]);
const categoryOptions = ref([]);

const {
  items,
  loading,
  saving,
  dialogVisible,
  currentItem,
  loadItems,
  openCreate,
  openEdit,
  saveItem,
  deleteItem,
} = useCrud('/products', () => ({ name: '', sku: '', brand_id: null, category_id: null }));

const dialogTitle = computed(() =>
  currentItem?.value?.id ? 'Редактирование продукта' : 'Новый продукт'
);

const findLabel = (options, value) => {
  const option = options.value.find((item) => item.id === value);
  return option ? option.name : value;
};

const loadOptions = async () => {
  const [brandsResponse, categoriesResponse] = await Promise.all([
    api.get('/brands'),
    api.get('/categories'),
  ]);
  brandOptions.value = brandsResponse.data;
  categoryOptions.value = categoriesResponse.data;
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

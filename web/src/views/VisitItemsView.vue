<template>
  <section class="page">
    <header class="flex align-items-center justify-content-between mb-4">
      <h2 class="m-0">Позиции визитов</h2>
      <Button label="Добавить" icon="pi pi-plus" @click="openCreate" />
    </header>
    <DataTable :value="items" dataKey="id" :loading="loading" responsiveLayout="scroll">
      <Column header="Визит">
        <template #body="{ data }">
          {{ visitLabelById[data.visit_id] ?? data.visit_id ?? '' }}
        </template>
      </Column>
      <Column header="Продукт">
        <template #body="{ data }">
          {{ productLabelById[data.product_id] ?? data.product_id ?? '' }}
        </template>
      </Column>
      <Column field="present_quantity" header="Кол-во на витрине" />
      <Column field="store_quantity" header="Кол-во на складе" />
      <Column field="price" header="Цена">
        <template #body="{ data }">
          {{ formatPrice(data.price) }}
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
        <div>
          <label class="block mb-2">Визит</label>
          <Dropdown
            v-model="currentItem.visit_id"
            :options="visitOptions"
            optionLabel="label"
            optionValue="id"
            placeholder="Выберите визит"
            class="w-full"
            required
          />
        </div>
        <div>
          <label class="block mb-2">Продукт</label>
          <Dropdown
            v-model="currentItem.product_id"
            :options="productOptions"
            optionLabel="name"
            optionValue="id"
            placeholder="Выберите продукт"
            class="w-full"
            required
          />
        </div>
        <div class="grid form-grid">
          <div class="col-12 md:col-4">
            <label class="block mb-2">На витрине</label>
            <InputNumber v-model="currentItem.present_quantity" class="w-full" />
          </div>
          <div class="col-12 md:col-4">
            <label class="block mb-2">На складе</label>
            <InputNumber v-model="currentItem.store_quantity" class="w-full" />
          </div>
          <div class="col-12 md:col-4">
            <label class="block mb-2">Цена</label>
            <InputNumber v-model="currentItem.price" mode="currency" currency="RUB" locale="ru-RU" class="w-full" />
          </div>
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

const visitOptions = ref([]);
const productOptions = ref([]);

const visitLabelById = computed(() =>
  Object.fromEntries(visitOptions.value.map((visit) => [visit.id, visit.label]))
);

const productLabelById = computed(() =>
  Object.fromEntries(productOptions.value.map((product) => [product.id, product.name]))
);

const crud = useCrud(
  '/visit-items',
  () => ({
    visit_id: null,
    product_id: null,
    present_quantity: null,
    store_quantity: null,
    price: null,
  }),
  {
    preparePayload: (payload) => ({
      ...payload,
      present_quantity: payload.present_quantity ?? null,
      store_quantity: payload.store_quantity ?? null,
      price:
        payload.price === null || payload.price === undefined || payload.price === ''
          ? null
          : Number(payload.price),
    }),
  }
);

const { items, loading, saving, dialogVisible, currentItem, loadItems, deleteItem } = crud;

const dialogTitle = computed(() =>
  currentItem?.value?.id ? 'Редактирование позиции' : 'Новая позиция'
);

const openCreate = () => {
  crud.openCreate();
};

const openEdit = (item) => {
  crud.openEdit({ ...item });
};

const saveItem = () => crud.saveItem();

const formatPrice = (value) => {
  if (value === null || value === undefined || value === '') {
    return '';
  }
  return Number(value).toFixed(2);
};

const loadOptions = async () => {
  const [visitsResponse, productsResponse] = await Promise.all([
    api.get('/visits'),
    api.get('/products'),
  ]);
  visitOptions.value = visitsResponse.data.map((visit) => ({
    id: visit.id,
    label: `${visit.id ? visit.id.slice(0, 6) : ''} · ${visit.visited_at ? new Date(visit.visited_at).toLocaleDateString('ru-RU') : ''}`,
  }));
  productOptions.value = productsResponse.data;
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

.form-grid {
  gap: 1rem;
}
</style>

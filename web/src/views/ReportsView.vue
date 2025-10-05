<template>
  <section class="page">
    <header class="flex align-items-center justify-content-between mb-4">
      <h2 class="m-0">Отчёт по визитам компании</h2>
    </header>
    <div class="grid">
      <div class="col-12 md:col-4">
        <label class="block mb-2">Компания</label>
        <Dropdown
          v-model="selectedCompany"
          :options="companyOptions"
          optionLabel="name"
          optionValue="id"
          placeholder="Выберите компанию"
          class="w-full"
          @change="fetchSummary"
        />
      </div>
    </div>

    <Card v-if="summary" class="mt-4">
      <template #title>
        Итоги для компании {{ resolveCompanyName(summary.company_id) }}
      </template>
      <template #content>
        <div class="grid text-lg">
          <div class="col-12 md:col-4">
            <span class="text-500 block mb-1">Всего визитов</span>
            <span class="font-bold">{{ summary.total_visits }}</span>
          </div>
          <div class="col-12 md:col-4">
            <span class="text-500 block mb-1">Продано единиц</span>
            <span class="font-bold">{{ summary.total_items }}</span>
          </div>
          <div class="col-12 md:col-4">
            <span class="text-500 block mb-1">Сумма продаж</span>
            <span class="font-bold">{{ summary.total_amount.toFixed(2) }} ₽</span>
          </div>
        </div>
      </template>
    </Card>

    <div v-else-if="selectedCompany && loading" class="mt-4">
      Загрузка отчёта...
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import api from '../services/api';

const companyOptions = ref([]);
const selectedCompany = ref(null);
const summary = ref(null);
const loading = ref(false);

const resolveCompanyName = (id) => {
  const company = companyOptions.value.find((option) => option.id === id);
  return company ? company.name : id;
};

const fetchSummary = async () => {
  summary.value = null;
  if (!selectedCompany.value) return;
  loading.value = true;
  try {
    const { data } = await api.get(`/reports/companies/${selectedCompany.value}/visits`);
    summary.value = data;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
};

onMounted(async () => {
  const { data } = await api.get('/companies');
  companyOptions.value = data;
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

"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import {
  PageHeader,
  DataTable,
  Button,
  Modal,
  Input,
  Select,
} from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Company, Worksite } from "@/lib/types";

interface WorksiteFormData {
  name: string;
  address: string;
  latitude: string;
  longitude: string;
  companyId: string;
}

const emptyForm: WorksiteFormData = {
  name: "",
  address: "",
  latitude: "",
  longitude: "",
  companyId: "",
};

export default function WorksitesPage() {
  const [worksites, setWorksites] = useState<Worksite[]>([]);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [editingWorksite, setEditingWorksite] = useState<Worksite | null>(null);
  const [formData, setFormData] = useState<WorksiteFormData>(emptyForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const companyLookup = useMemo(() => {
    const map: Record<string, string> = {};
    for (const c of companies) {
      map[c.id] = c.name;
    }
    return map;
  }, [companies]);

  const companyOptions = useMemo(
    () => companies.map((c) => ({ value: c.id, label: c.name })),
    [companies]
  );

  const fetchData = useCallback(() => {
    setLoading(true);
    setError(null);
    Promise.all([
      apiClient<Worksite[]>("/api/v1/worksites"),
      apiClient<Company[]>("/api/v1/companies"),
    ])
      .then(([ws, cs]) => {
        setWorksites(ws);
        setCompanies(cs);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  function openCreateModal() {
    setEditingWorksite(null);
    setFormData(emptyForm);
    setFormError(null);
    setModalOpen(true);
  }

  function openEditModal(worksite: Worksite) {
    setEditingWorksite(worksite);
    setFormData({
      name: worksite.name,
      address: worksite.address ?? "",
      latitude: worksite.latitude != null ? String(worksite.latitude) : "",
      longitude: worksite.longitude != null ? String(worksite.longitude) : "",
      companyId: worksite.companyId,
    });
    setFormError(null);
    setModalOpen(true);
  }

  function closeModal() {
    setModalOpen(false);
    setEditingWorksite(null);
    setFormData(emptyForm);
    setFormError(null);
  }

  function handleChange(field: keyof WorksiteFormData, value: string) {
    setFormData((prev) => ({ ...prev, [field]: value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!formData.name.trim()) {
      setFormError("Worksite name is required.");
      return;
    }

    if (!formData.companyId) {
      setFormError("Please select a company.");
      return;
    }

    setSubmitting(true);
    setFormError(null);

    const body: Record<string, unknown> = {
      name: formData.name.trim(),
      companyId: formData.companyId,
      address: formData.address.trim() || undefined,
    };

    if (formData.latitude.trim()) {
      const lat = parseFloat(formData.latitude.trim());
      if (isNaN(lat)) {
        setFormError("Latitude must be a valid number.");
        setSubmitting(false);
        return;
      }
      body.latitude = lat;
    }

    if (formData.longitude.trim()) {
      const lng = parseFloat(formData.longitude.trim());
      if (isNaN(lng)) {
        setFormError("Longitude must be a valid number.");
        setSubmitting(false);
        return;
      }
      body.longitude = lng;
    }

    try {
      if (editingWorksite) {
        await apiClient<Worksite>(
          `/api/v1/worksites/${editingWorksite.id}`,
          { method: "PUT", body: JSON.stringify(body) }
        );
      } else {
        await apiClient<Worksite>("/api/v1/worksites", {
          method: "POST",
          body: JSON.stringify(body),
        });
      }
      closeModal();
      fetchData();
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "An error occurred.");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return <div className="text-sm text-gray-500">Loading worksites...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load worksites: {error}
      </div>
    );
  }

  const columns = [
    { key: "name", header: "Name" },
    {
      key: "companyId",
      header: "Company",
      render: (w: Worksite) => companyLookup[w.companyId] ?? "\u2014",
    },
    { key: "address", header: "Address" },
    {
      key: "latitude",
      header: "Coordinates",
      render: (w: Worksite) =>
        w.latitude != null && w.longitude != null
          ? `${w.latitude.toFixed(4)}, ${w.longitude.toFixed(4)}`
          : "\u2014",
    },
  ];

  return (
    <div>
      <PageHeader
        title="Worksites"
        description="View and manage worksites across all companies."
        action={<Button onClick={openCreateModal}>Add Worksite</Button>}
      />
      <DataTable
        columns={columns}
        data={worksites}
        emptyMessage="No worksites found."
        onRowClick={openEditModal}
      />

      <Modal
        open={modalOpen}
        onClose={closeModal}
        title={editingWorksite ? "Edit Worksite" : "Add Worksite"}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Worksite Name"
            required
            value={formData.name}
            onChange={(e) => handleChange("name", e.target.value)}
            placeholder="e.g. Manchester Arena"
          />
          <Select
            label="Company"
            required
            options={companyOptions}
            value={formData.companyId}
            onChange={(e) => handleChange("companyId", e.target.value)}
          />
          <Input
            label="Address"
            value={formData.address}
            onChange={(e) => handleChange("address", e.target.value)}
            placeholder="e.g. Victoria Station Approach, Manchester"
          />
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Latitude"
              type="text"
              inputMode="decimal"
              value={formData.latitude}
              onChange={(e) => handleChange("latitude", e.target.value)}
              placeholder="e.g. 53.4884"
            />
            <Input
              label="Longitude"
              type="text"
              inputMode="decimal"
              value={formData.longitude}
              onChange={(e) => handleChange("longitude", e.target.value)}
              placeholder="e.g. -2.2446"
            />
          </div>

          {formError && (
            <p className="text-sm text-red-600">{formError}</p>
          )}

          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="secondary" onClick={closeModal}>
              Cancel
            </Button>
            <Button type="submit" disabled={submitting}>
              {submitting
                ? "Saving..."
                : editingWorksite
                  ? "Update Worksite"
                  : "Create Worksite"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

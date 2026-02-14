"use client";

import { useCallback, useEffect, useState } from "react";
import { PageHeader, DataTable, Button, Modal, Input } from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Company } from "@/lib/types";

interface CompanyFormData {
  name: string;
  address: string;
  phone: string;
  email: string;
}

const emptyForm: CompanyFormData = {
  name: "",
  address: "",
  phone: "",
  email: "",
};

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [editingCompany, setEditingCompany] = useState<Company | null>(null);
  const [formData, setFormData] = useState<CompanyFormData>(emptyForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const fetchCompanies = useCallback(() => {
    setLoading(true);
    setError(null);
    apiClient<Company[]>("/api/v1/companies")
      .then(setCompanies)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    fetchCompanies();
  }, [fetchCompanies]);

  function openCreateModal() {
    setEditingCompany(null);
    setFormData(emptyForm);
    setFormError(null);
    setModalOpen(true);
  }

  function openEditModal(company: Company) {
    setEditingCompany(company);
    setFormData({
      name: company.name,
      address: company.address ?? "",
      phone: company.phone ?? "",
      email: company.email ?? "",
    });
    setFormError(null);
    setModalOpen(true);
  }

  function closeModal() {
    setModalOpen(false);
    setEditingCompany(null);
    setFormData(emptyForm);
    setFormError(null);
  }

  function handleChange(field: keyof CompanyFormData, value: string) {
    setFormData((prev) => ({ ...prev, [field]: value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!formData.name.trim()) {
      setFormError("Company name is required.");
      return;
    }

    setSubmitting(true);
    setFormError(null);

    const body = {
      name: formData.name.trim(),
      address: formData.address.trim() || undefined,
      phone: formData.phone.trim() || undefined,
      email: formData.email.trim() || undefined,
    };

    try {
      if (editingCompany) {
        await apiClient<Company>(`/api/v1/companies/${editingCompany.id}`, {
          method: "PUT",
          body: JSON.stringify(body),
        });
      } else {
        await apiClient<Company>("/api/v1/companies", {
          method: "POST",
          body: JSON.stringify(body),
        });
      }
      closeModal();
      fetchCompanies();
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "An error occurred.");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return <div className="text-sm text-gray-500">Loading companies...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load companies: {error}
      </div>
    );
  }

  const columns = [
    { key: "name", header: "Name" },
    { key: "email", header: "Email" },
    { key: "phone", header: "Telephone" },
    {
      key: "address",
      header: "Address",
      render: (c: Company) => (
        <span className="max-w-xs truncate">{c.address ?? "\u2014"}</span>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Companies"
        description="Manage security companies registered on the platform."
        action={<Button onClick={openCreateModal}>Add Company</Button>}
      />
      <DataTable
        columns={columns}
        data={companies}
        emptyMessage="No companies found. Add one to get started."
        onRowClick={openEditModal}
      />

      <Modal
        open={modalOpen}
        onClose={closeModal}
        title={editingCompany ? "Edit Company" : "Add Company"}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Company Name"
            required
            value={formData.name}
            onChange={(e) => handleChange("name", e.target.value)}
            placeholder="e.g. Sentinel Security Ltd"
          />
          <Input
            label="Address"
            value={formData.address}
            onChange={(e) => handleChange("address", e.target.value)}
            placeholder="e.g. 42 High Street, London"
          />
          <Input
            label="Telephone"
            type="tel"
            value={formData.phone}
            onChange={(e) => handleChange("phone", e.target.value)}
            placeholder="e.g. 020 7946 0958"
          />
          <Input
            label="Email"
            type="email"
            value={formData.email}
            onChange={(e) => handleChange("email", e.target.value)}
            placeholder="e.g. info@sentinelsecurity.co.uk"
          />

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
                : editingCompany
                  ? "Update Company"
                  : "Create Company"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

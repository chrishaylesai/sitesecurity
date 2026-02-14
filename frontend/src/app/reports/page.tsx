"use client";

import { useEffect, useState, useCallback } from "react";
import {
  PageHeader,
  DataTable,
  Button,
  Modal,
  Input,
  TextArea,
  Select,
} from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { ShiftReportTemplate, ShiftReport, Company } from "@/lib/types";

export default function ReportsPage() {
  const [templates, setTemplates] = useState<ShiftReportTemplate[]>([]);
  const [reports, setReports] = useState<ShiftReport[]>([]);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [form, setForm] = useState({
    name: "",
    companyId: "",
    fields: "",
  });

  const companyLookup = companies.reduce<Record<string, string>>(
    (acc, c) => {
      acc[c.id] = c.name;
      return acc;
    },
    {}
  );

  const templateLookup = templates.reduce<Record<string, string>>(
    (acc, t) => {
      acc[t.id] = t.name;
      return acc;
    },
    {}
  );

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [templatesData, reportsData, companiesData] = await Promise.all([
        apiClient<ShiftReportTemplate[]>("/api/v1/shift-reports/templates"),
        apiClient<ShiftReport[]>("/api/v1/shift-reports"),
        apiClient<Company[]>("/api/v1/companies"),
      ]);
      setTemplates(templatesData);
      setReports(reportsData);
      setCompanies(companiesData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleCreate = async () => {
    if (!form.name.trim()) return;
    setSubmitting(true);
    try {
      await apiClient<ShiftReportTemplate>("/api/v1/shift-reports/templates", {
        method: "POST",
        body: JSON.stringify({
          name: form.name,
          companyId: form.companyId || undefined,
          fields: form.fields || undefined,
        }),
      });
      setModalOpen(false);
      setForm({ name: "", companyId: "", fields: "" });
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create template");
    } finally {
      setSubmitting(false);
    }
  };

  const formatDate = (iso: string) => {
    const d = new Date(iso);
    return d.toLocaleDateString("en-GB", {
      day: "numeric",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (loading) {
    return <div className="text-sm text-gray-500">Loading reports...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load reports: {error}
      </div>
    );
  }

  const templateColumns = [
    { key: "name", header: "Name" },
    {
      key: "companyId",
      header: "Company",
      render: (t: ShiftReportTemplate) => companyLookup[t.companyId] || "\u2014",
    },
    {
      key: "createdAt",
      header: "Created At",
      render: (t: ShiftReportTemplate) => formatDate(t.createdAt),
    },
  ];

  const reportColumns = [
    { key: "shiftId", header: "Shift ID" },
    { key: "workerId", header: "Worker ID" },
    {
      key: "submittedAt",
      header: "Submitted At",
      render: (r: ShiftReport) => formatDate(r.submittedAt),
    },
    {
      key: "templateId",
      header: "Template",
      render: (r: ShiftReport) =>
        r.templateId ? templateLookup[r.templateId] || "\u2014" : "\u2014",
    },
  ];

  const companyOptions = companies.map((c) => ({
    value: c.id,
    label: c.name,
  }));

  return (
    <div>
      <PageHeader
        title="Shift Reports"
        description="Manage report templates and view submitted shift reports."
        action={
          <Button onClick={() => setModalOpen(true)}>Create Template</Button>
        }
      />

      <div className="space-y-8">
        <section>
          <h2 className="mb-4 text-lg font-semibold text-gray-900">
            Report Templates
          </h2>
          <DataTable
            columns={templateColumns}
            data={templates}
            emptyMessage="No report templates found."
          />
        </section>

        <section>
          <h2 className="mb-4 text-lg font-semibold text-gray-900">
            Submitted Reports
          </h2>
          <DataTable
            columns={reportColumns}
            data={reports}
            emptyMessage="No submitted reports found."
          />
        </section>
      </div>

      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title="Create Template"
      >
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleCreate();
          }}
          className="space-y-4"
        >
          <Input
            label="Name"
            required
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
          />
          <Select
            label="Company"
            options={companyOptions}
            value={form.companyId}
            onChange={(e) => setForm({ ...form, companyId: e.target.value })}
          />
          <TextArea
            label="Fields"
            placeholder='[{"label": "Notes", "type": "text", "required": true}]'
            value={form.fields}
            onChange={(e) => setForm({ ...form, fields: e.target.value })}
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button
              type="button"
              variant="secondary"
              onClick={() => setModalOpen(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={submitting || !form.name.trim()}>
              {submitting ? "Creating..." : "Create Template"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

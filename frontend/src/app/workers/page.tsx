"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { PageHeader, DataTable, Button, Modal, Input } from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Worker } from "@/lib/types";

export default function WorkersPage() {
  const router = useRouter();
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [form, setForm] = useState({
    firstName: "",
    lastName: "",
    email: "",
    phone: "",
    authSubject: "",
  });

  const fetchWorkers = () => {
    setLoading(true);
    apiClient<Worker[]>("/api/v1/workers")
      .then(setWorkers)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchWorkers();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await apiClient<Worker>("/api/v1/workers", {
        method: "POST",
        body: JSON.stringify({
          firstName: form.firstName,
          lastName: form.lastName,
          email: form.email,
          phone: form.phone || undefined,
          authSubject: form.authSubject,
        }),
      });
      setModalOpen(false);
      setForm({ firstName: "", lastName: "", email: "", phone: "", authSubject: "" });
      fetchWorkers();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create worker");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="text-sm text-gray-500">Loading workers...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load workers: {error}
      </div>
    );
  }

  const columns = [
    {
      key: "name",
      header: "Name",
      render: (w: Worker) => `${w.firstName} ${w.lastName}`,
    },
    { key: "email", header: "Email" },
    {
      key: "phone",
      header: "Telephone",
      render: (w: Worker) => w.phone ?? "\u2014",
    },
  ];

  return (
    <div>
      <PageHeader
        title="Workers"
        description="Manage security workers, their qualifications, and company memberships."
        action={<Button onClick={() => setModalOpen(true)}>Add Worker</Button>}
      />
      <DataTable
        columns={columns}
        data={workers}
        emptyMessage="No workers found."
        onRowClick={(w) => router.push(`/workers/${w.id}`)}
      />

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title="Add Worker">
        <form onSubmit={handleCreate} className="space-y-4">
          <Input
            label="First Name"
            value={form.firstName}
            onChange={(e) => setForm({ ...form, firstName: e.target.value })}
            required
          />
          <Input
            label="Last Name"
            value={form.lastName}
            onChange={(e) => setForm({ ...form, lastName: e.target.value })}
            required
          />
          <Input
            label="Email"
            type="email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
          />
          <Input
            label="Telephone"
            type="tel"
            value={form.phone}
            onChange={(e) => setForm({ ...form, phone: e.target.value })}
          />
          <Input
            label="Auth Subject"
            value={form.authSubject}
            onChange={(e) => setForm({ ...form, authSubject: e.target.value })}
            required
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button
              type="button"
              variant="secondary"
              onClick={() => setModalOpen(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={submitting}>
              {submitting ? "Creating..." : "Create Worker"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

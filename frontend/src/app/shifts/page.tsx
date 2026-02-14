"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import {
  PageHeader,
  DataTable,
  Badge,
  Button,
  Modal,
  Input,
  TextArea,
  Select,
} from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Shift, Worksite } from "@/lib/types";

const statusVariant: Record<string, string> = {
  open: "info",
  assigned: "default",
  in_progress: "warning",
  completed: "success",
  cancelled: "danger",
};

export default function ShiftsPage() {
  const router = useRouter();
  const [shifts, setShifts] = useState<Shift[]>([]);
  const [worksites, setWorksites] = useState<Worksite[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [form, setForm] = useState({
    title: "",
    description: "",
    worksiteId: "",
    startTime: "",
    endTime: "",
  });

  const worksiteLookup = worksites.reduce<Record<string, string>>(
    (acc, ws) => {
      acc[ws.id] = ws.name;
      return acc;
    },
    {}
  );

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [shiftsData, worksitesData] = await Promise.all([
        apiClient<Shift[]>("/api/v1/shifts"),
        apiClient<Worksite[]>("/api/v1/worksites"),
      ]);
      setShifts(shiftsData);
      setWorksites(worksitesData);
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
    if (!form.title.trim()) return;
    setSubmitting(true);
    try {
      await apiClient<Shift>("/api/v1/shifts", {
        method: "POST",
        body: JSON.stringify({
          title: form.title,
          description: form.description || undefined,
          worksiteId: form.worksiteId || undefined,
          startTime: form.startTime ? new Date(form.startTime).toISOString() : undefined,
          endTime: form.endTime ? new Date(form.endTime).toISOString() : undefined,
        }),
      });
      setModalOpen(false);
      setForm({ title: "", description: "", worksiteId: "", startTime: "", endTime: "" });
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create shift");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="text-sm text-gray-500">Loading shifts...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load shifts: {error}
      </div>
    );
  }

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

  const columns = [
    { key: "title", header: "Title" },
    {
      key: "worksiteId",
      header: "Worksite",
      render: (s: Shift) => worksiteLookup[s.worksiteId] || "\u2014",
    },
    {
      key: "startTime",
      header: "Start",
      render: (s: Shift) => formatDate(s.startTime),
    },
    {
      key: "endTime",
      header: "End",
      render: (s: Shift) => formatDate(s.endTime),
    },
    {
      key: "status",
      header: "Status",
      render: (s: Shift) => (
        <Badge
          variant={
            statusVariant[s.status] as
              | "info"
              | "default"
              | "warning"
              | "success"
              | "danger"
          }
        >
          {s.status.replace("_", " ")}
        </Badge>
      ),
    },
  ];

  const worksiteOptions = worksites.map((ws) => ({
    value: ws.id,
    label: ws.name,
  }));

  return (
    <div>
      <PageHeader
        title="Shifts"
        description="Schedule, assign, and track security shifts."
        action={<Button onClick={() => setModalOpen(true)}>Create Shift</Button>}
      />
      <DataTable
        columns={columns}
        data={shifts}
        emptyMessage="No shifts scheduled."
        onRowClick={(s: Shift) => router.push(`/shifts/${s.id}`)}
      />

      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title="Create Shift"
      >
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleCreate();
          }}
          className="space-y-4"
        >
          <Input
            label="Title"
            required
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
          />
          <TextArea
            label="Description"
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
          />
          <Select
            label="Worksite"
            options={worksiteOptions}
            value={form.worksiteId}
            onChange={(e) => setForm({ ...form, worksiteId: e.target.value })}
          />
          <Input
            label="Start Time"
            type="datetime-local"
            value={form.startTime}
            onChange={(e) => setForm({ ...form, startTime: e.target.value })}
          />
          <Input
            label="End Time"
            type="datetime-local"
            value={form.endTime}
            onChange={(e) => setForm({ ...form, endTime: e.target.value })}
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button
              type="button"
              variant="secondary"
              onClick={() => setModalOpen(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={submitting || !form.title.trim()}>
              {submitting ? "Creating..." : "Create Shift"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

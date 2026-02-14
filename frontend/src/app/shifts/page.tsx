"use client";

import { useEffect, useState } from "react";
import { PageHeader, DataTable, Badge, Button } from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Shift } from "@/lib/types";

const statusVariant: Record<string, string> = {
  open: "info",
  assigned: "default",
  in_progress: "warning",
  completed: "success",
  cancelled: "danger",
};

export default function ShiftsPage() {
  const [shifts, setShifts] = useState<Shift[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiClient<Shift[]>("/api/v1/shifts")
      .then(setShifts)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

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
        <Badge variant={statusVariant[s.status] as "info" | "default" | "warning" | "success" | "danger"}>
          {s.status.replace("_", " ")}
        </Badge>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Shifts"
        description="Schedule, assign, and track security shifts."
        action={<Button>Create Shift</Button>}
      />
      <DataTable
        columns={columns}
        data={shifts}
        emptyMessage="No shifts scheduled."
      />
    </div>
  );
}

"use client";

import { useEffect, useState } from "react";
import { PageHeader, DataTable, Badge } from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Alarm } from "@/lib/types";

const statusVariant: Record<string, string> = {
  raised: "danger",
  acknowledged: "warning",
  resolved: "success",
};

export default function AlarmsPage() {
  const [alarms, setAlarms] = useState<Alarm[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiClient<Alarm[]>("/api/v1/alarms")
      .then(setAlarms)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="text-sm text-gray-500">Loading alarms...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load alarms: {error}
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
    {
      key: "status",
      header: "Status",
      render: (a: Alarm) => (
        <Badge variant={statusVariant[a.status] as "danger" | "warning" | "success"}>
          {a.status}
        </Badge>
      ),
    },
    {
      key: "message",
      header: "Message",
      render: (a: Alarm) => a.message ?? "No message",
    },
    {
      key: "raisedAt",
      header: "Raised",
      render: (a: Alarm) => formatDate(a.raisedAt),
    },
    {
      key: "location",
      header: "Location",
      render: (a: Alarm) =>
        a.latitude && a.longitude
          ? `${a.latitude.toFixed(4)}, ${a.longitude.toFixed(4)}`
          : "—",
    },
  ];

  return (
    <div>
      <PageHeader
        title="Alarms"
        description="Monitor and respond to alarms raised by security workers."
      />
      <DataTable
        columns={columns}
        data={alarms}
        emptyMessage="No alarms to display."
      />
    </div>
  );
}

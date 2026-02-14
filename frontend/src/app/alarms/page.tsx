"use client";

import { useCallback, useEffect, useState } from "react";
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
import type { Alarm, Worker } from "@/lib/types";

const statusVariant: Record<string, "danger" | "warning" | "success"> = {
  raised: "danger",
  acknowledged: "warning",
  resolved: "success",
};

export default function AlarmsPage() {
  const [alarms, setAlarms] = useState<Alarm[]>([]);
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [form, setForm] = useState({
    message: "",
    latitude: "",
    longitude: "",
    workerId: "",
    shiftId: "",
  });

  const workerMap = new Map(workers.map((w) => [w.id, w]));

  const fetchData = useCallback(async () => {
    try {
      const [alarmsData, workersData] = await Promise.all([
        apiClient<Alarm[]>("/api/v1/alarms"),
        apiClient<Worker[]>("/api/v1/workers"),
      ]);
      setAlarms(alarmsData);
      setWorkers(workersData);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleRaiseAlarm = async () => {
    setSubmitting(true);
    try {
      await apiClient("/api/v1/alarms", {
        method: "POST",
        body: JSON.stringify({
          workerId: form.workerId,
          message: form.message || undefined,
          latitude: form.latitude ? parseFloat(form.latitude) : undefined,
          longitude: form.longitude ? parseFloat(form.longitude) : undefined,
          shiftId: form.shiftId || undefined,
        }),
      });
      setModalOpen(false);
      setForm({ message: "", latitude: "", longitude: "", workerId: "", shiftId: "" });
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to raise alarm");
    } finally {
      setSubmitting(false);
    }
  };

  const handleAcknowledge = async (id: string) => {
    try {
      await apiClient(`/api/v1/alarms/${id}/acknowledge`, { method: "PATCH" });
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to acknowledge alarm");
    }
  };

  const handleResolve = async (id: string) => {
    try {
      await apiClient(`/api/v1/alarms/${id}/resolve`, { method: "PATCH" });
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to resolve alarm");
    }
  };

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

  const workerName = (workerId: string) => {
    const w = workerMap.get(workerId);
    return w ? `${w.firstName} ${w.lastName}` : workerId;
  };

  const workerOptions = workers.map((w) => ({
    value: w.id,
    label: `${w.firstName} ${w.lastName}`,
  }));

  const columns = [
    {
      key: "status",
      header: "Status",
      render: (a: Alarm) => (
        <Badge variant={statusVariant[a.status]}>{a.status}</Badge>
      ),
    },
    {
      key: "worker",
      header: "Worker",
      render: (a: Alarm) => workerName(a.workerId),
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
        a.latitude != null && a.longitude != null
          ? `${a.latitude.toFixed(4)}, ${a.longitude.toFixed(4)}`
          : "\u2014",
    },
    {
      key: "actions",
      header: "Actions",
      render: (a: Alarm) => {
        if (a.status === "raised") {
          return (
            <Button
              variant="secondary"
              size="sm"
              onClick={() => handleAcknowledge(a.id)}
            >
              Acknowledge
            </Button>
          );
        }
        if (a.status === "acknowledged") {
          return (
            <Button
              variant="primary"
              size="sm"
              onClick={() => handleResolve(a.id)}
            >
              Resolve
            </Button>
          );
        }
        return (
          <span className="text-green-600" title="Resolved">
            &#10003;
          </span>
        );
      },
    },
  ];

  return (
    <div>
      <PageHeader
        title="Alarms"
        description="Monitor and respond to alarms raised by security workers."
      />
      <div className="mb-4 flex justify-end">
        <Button variant="danger" onClick={() => setModalOpen(true)}>
          Raise Alarm
        </Button>
      </div>
      <DataTable
        columns={columns}
        data={alarms}
        emptyMessage="No alarms to display."
      />
      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title="Raise Alarm"
      >
        <div className="space-y-4">
          <Select
            label="Worker"
            options={workerOptions}
            value={form.workerId}
            onChange={(e) => setForm({ ...form, workerId: e.target.value })}
          />
          <TextArea
            label="Message (optional)"
            value={form.message}
            onChange={(e) => setForm({ ...form, message: e.target.value })}
            placeholder="Describe the alarm..."
          />
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Latitude"
              type="number"
              step="any"
              value={form.latitude}
              onChange={(e) => setForm({ ...form, latitude: e.target.value })}
              placeholder="e.g. 51.5074"
            />
            <Input
              label="Longitude"
              type="number"
              step="any"
              value={form.longitude}
              onChange={(e) => setForm({ ...form, longitude: e.target.value })}
              placeholder="e.g. -0.1278"
            />
          </div>
          <Input
            label="Shift ID (optional)"
            value={form.shiftId}
            onChange={(e) => setForm({ ...form, shiftId: e.target.value })}
            placeholder="Enter shift ID"
          />
          <div className="flex justify-end gap-2 pt-2">
            <Button
              variant="secondary"
              onClick={() => setModalOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="danger"
              onClick={handleRaiseAlarm}
              disabled={!form.workerId || submitting}
            >
              {submitting ? "Raising..." : "Raise Alarm"}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

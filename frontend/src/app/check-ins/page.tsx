"use client";

import { useCallback, useEffect, useState } from "react";
import {
  PageHeader,
  DataTable,
  Button,
  Modal,
  Input,
  Select,
} from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { LocationCheckIn, Worker } from "@/lib/types";

export default function CheckInsPage() {
  const [checkIns, setCheckIns] = useState<LocationCheckIn[]>([]);
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [form, setForm] = useState({
    workerId: "",
    latitude: "",
    longitude: "",
    shiftId: "",
  });

  const workerMap = new Map(workers.map((w) => [w.id, w]));

  const fetchData = useCallback(async () => {
    try {
      const [checkInsData, workersData] = await Promise.all([
        apiClient<LocationCheckIn[]>("/api/v1/check-ins"),
        apiClient<Worker[]>("/api/v1/workers"),
      ]);
      setCheckIns(checkInsData);
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

  const handleCreate = async () => {
    setSubmitting(true);
    try {
      await apiClient("/api/v1/check-ins", {
        method: "POST",
        body: JSON.stringify({
          workerId: form.workerId,
          latitude: parseFloat(form.latitude),
          longitude: parseFloat(form.longitude),
          shiftId: form.shiftId || undefined,
        }),
      });
      setModalOpen(false);
      setForm({ workerId: "", latitude: "", longitude: "", shiftId: "" });
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to record check-in");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="text-sm text-gray-500">Loading check-ins...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load check-ins: {error}
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
      key: "worker",
      header: "Worker",
      render: (c: LocationCheckIn) => workerName(c.workerId),
    },
    {
      key: "coordinates",
      header: "Coordinates",
      render: (c: LocationCheckIn) =>
        `${c.latitude.toFixed(4)}, ${c.longitude.toFixed(4)}`,
    },
    {
      key: "shiftId",
      header: "Shift ID",
      render: (c: LocationCheckIn) => c.shiftId ?? "\u2014",
    },
    {
      key: "recordedAt",
      header: "Recorded At",
      render: (c: LocationCheckIn) => formatDate(c.recordedAt),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Location Check-ins"
        description="View GPS check-in records from security workers."
      />
      <div className="mb-4 flex justify-end">
        <Button variant="primary" onClick={() => setModalOpen(true)}>
          Record Check-in
        </Button>
      </div>
      <DataTable
        columns={columns}
        data={checkIns}
        emptyMessage="No check-ins to display."
      />
      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title="Record Check-in"
      >
        <div className="space-y-4">
          <Select
            label="Worker"
            options={workerOptions}
            value={form.workerId}
            onChange={(e) => setForm({ ...form, workerId: e.target.value })}
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
              variant="primary"
              onClick={handleCreate}
              disabled={!form.workerId || !form.latitude || !form.longitude || submitting}
            >
              {submitting ? "Recording..." : "Record Check-in"}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

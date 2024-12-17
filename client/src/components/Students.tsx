import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { FaEdit, FaTrash } from "react-icons/fa";
import { BASE_URL } from "../App";

// Define Student type matching the backend struct
export type Student = {
    _id: string;
    firstName: string;
    lastName: string;
    phoneNumber: string;
    email: string;
    address: string;
};

const Students = () => {
    const { data: students, isLoading } = useQuery<Student[]>({
        queryKey: ["students"],
        queryFn: async () => {
            try {
                const res = await fetch(BASE_URL + "/students");
                if (!res.ok) {
                    // Throw an error if the response is not OK
                    throw new Error("Failed to fetch students");
                }
                const data = await res.json();
                return data || []; // Ensure we return an array
            } catch (error) {
                console.error("Error fetching students:", error);
                return []; // Return an empty array on error
            }
        },
        // Add retry option to prevent infinite retries
        retry: 1
    });

    // Placeholder functions for update and delete
    const handleUpdate = (student: Student) => {
        console.log("Update student:", student);
        // TODO: Implement update logic
    };

    const handleDelete = (student: Student) => {
        console.log("Delete student:", student);
        // TODO: Implement delete logic
    };

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4 text-center">Students List</h1>
            
            {isLoading && (
                <div className="text-center">
                    Loading...
                </div>
            )}
            
            {!isLoading && students?.length === 0 && (
                <div className="text-center">
                    <p>No students found</p>
                </div>
            )}
            
            {!isLoading && students && students.length > 0 && (
                <div className="overflow-x-auto">
                    <table className="w-full border-collapse border border-gray-200">
                        <thead className="bg-gray-100">
                            <tr>
                                <th className="border p-2 text-left">First Name</th>
                                <th className="border p-2 text-left">Last Name</th>
                                <th className="border p-2 text-left">Email</th>
                                <th className="border p-2 text-left">Phone Number</th>
                                <th className="border p-2 text-left">Address</th>
                                <th className="border p-2 text-center">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {students.map((student) => (
                                <tr key={student._id} className="hover:bg-gray-50 transition">
                                    <td className="border p-2">{student.firstName}</td>
                                    <td className="border p-2">{student.lastName}</td>
                                    <td className="border p-2">{student.email}</td>
                                    <td className="border p-2">{student.phoneNumber}</td>
                                    <td className="border p-2">{student.address}</td>
                                    <td className="border p-2">
                                        <div className="flex justify-center space-x-2">
                                            <button 
                                                onClick={() => handleUpdate(student)}
                                                className="text-blue-500 hover:text-blue-700 transition"
                                                title="Update Student"
                                            >
                                                <FaEdit size={20} />
                                            </button>
                                            <button 
                                                onClick={() => handleDelete(student)}
                                                className="text-red-500 hover:text-red-700 transition"
                                                title="Delete Student"
                                            >
                                                <FaTrash size={20} />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};

export default Students;